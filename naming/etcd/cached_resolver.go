package etcd

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	clientv3 "go.etcd.io/etcd/client/v3"
	gresolver "google.golang.org/grpc/resolver"
)

const (
	defaultCacheTTL            = 5 * time.Minute
	defaultHealthCheckInterval = 30 * time.Second
	defaultFallbackThreshold   = 3
)

type cachedAddrEntry struct {
	addrs     []gresolver.Address
	updatedAt time.Time
}

type CachedResolver struct {
	gresolver.Builder
	conf *naming.Config
	cli  *clientv3.Client

	cache map[string]*cachedAddrEntry
	mu    sync.RWMutex

	failCount atomic.Uint32
	fallback  atomic.Bool
	done      chan struct{}
	once      sync.Once
}

func NewCachedResolver(inner naming.Resolver, cli *clientv3.Client, conf *naming.Config) *CachedResolver {
	cr := &CachedResolver{
		Builder: inner,
		conf:    conf,
		cli:     cli,
		cache:   make(map[string]*cachedAddrEntry),
		done:    make(chan struct{}),
	}
	if conf.ResolverFallback && conf.ResolverCacheTTL > 0 {
		go cr.healthCheckLoop()
	}
	return cr
}

func (c *CachedResolver) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	wrappedCC := &cachedClientConn{
		ClientConn: cc,
		cached:     c,
		target:     target.Endpoint(),
	}
	return c.Builder.Build(target, wrappedCC, opts)
}

func (c *CachedResolver) Config() *naming.Config {
	return c.conf
}

func (c *CachedResolver) Close() error {
	c.once.Do(func() {
		close(c.done)
	})
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

func (c *CachedResolver) cacheTTL() time.Duration {
	if c.conf.ResolverCacheTTL > 0 {
		return c.conf.ResolverCacheTTL
	}
	return defaultCacheTTL
}

func (c *CachedResolver) healthCheckInterval() time.Duration {
	if c.conf.HealthCheckInterval > 0 {
		return c.conf.HealthCheckInterval
	}
	return defaultHealthCheckInterval
}

func (c *CachedResolver) fallbackThreshold() int {
	if c.conf.FallbackThreshold > 0 {
		return c.conf.FallbackThreshold
	}
	return defaultFallbackThreshold
}

func (c *CachedResolver) updateCache(target string, addrs []gresolver.Address) {
	c.mu.Lock()
	c.cache[target] = &cachedAddrEntry{
		addrs:     copyAddrs(addrs),
		updatedAt: time.Now(),
	}
	c.mu.Unlock()

	c.failCount.Store(0)
	if c.fallback.Swap(false) {
		logger.DefaultLogger.Infof("[resolver] restored from etcd: %s", target)
	}
}

func (c *CachedResolver) getCache(target string) ([]gresolver.Address, bool) {
	c.mu.RLock()
	entry, ok := c.cache[target]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Since(entry.updatedAt) > c.cacheTTL() {
		c.mu.Lock()
		delete(c.cache, target)
		c.mu.Unlock()
		logger.DefaultLogger.Warnf("[resolver] cache expired: %s", target)
		return nil, false
	}
	return copyAddrs(entry.addrs), true
}

func (c *CachedResolver) markFailure(target string) {
	count := c.failCount.Add(1)
	if uint32(c.fallbackThreshold()) > 0 && count >= uint32(c.fallbackThreshold()) {
		if c.fallback.CompareAndSwap(false, true) {
			logger.DefaultLogger.Warnf("[resolver] fallback to cached: %s (failures: %d)", target, count)
		}
	}
}

func (c *CachedResolver) isFallback() bool {
	return c.fallback.Load()
}

func (c *CachedResolver) healthCheckLoop() {
	ticker := time.NewTicker(c.healthCheckInterval())
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.checkEtcdHealth()
		}
	}
}

func (c *CachedResolver) checkEtcdHealth() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.cli.Get(ctx, "health_check", clientv3.WithSerializable())
	if err != nil {
		c.failCount.Add(1)
	} else {
		c.failCount.Store(0)
	}
}

type cachedClientConn struct {
	gresolver.ClientConn
	cached *CachedResolver
	target string
}

func (c *cachedClientConn) UpdateState(state gresolver.State) error {
	c.cached.updateCache(c.target, state.Addresses)
	if c.cached.isFallback() {
		return nil
	}
	return c.ClientConn.UpdateState(state)
}

func (c *cachedClientConn) ReportError(err error) {
	c.cached.markFailure(c.target)
	c.ClientConn.ReportError(err)
}

func copyAddrs(addrs []gresolver.Address) []gresolver.Address {
	out := make([]gresolver.Address, len(addrs))
	copy(out, addrs)
	return out
}
