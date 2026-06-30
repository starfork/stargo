package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/starfork/stargo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64
}

type RateLimiter struct {
	limiters  sync.Map // map[key]*limiterEntry
	rate      rate.Limit
	burst     int
	cleanUpIn time.Duration
	stopChan  chan struct{}
}

func NewRateLimiter(r rate.Limit, burst int, cleanUpIn time.Duration) *RateLimiter {
	krl := &RateLimiter{
		rate:      r,
		burst:     burst,
		cleanUpIn: cleanUpIn,
		stopChan:  make(chan struct{}),
	}
	if cleanUpIn > 0 {
		go krl.cleanup()
	}
	return krl
}

// getLimiter 获取或创建一个限流器
func (k *RateLimiter) getLimiter(key string) *rate.Limiter {
	entry, _ := k.limiters.LoadOrStore(key, &limiterEntry{
		limiter: rate.NewLimiter(k.rate, k.burst),
	})
	e := entry.(*limiterEntry)
	e.lastSeen.Store(time.Now().UnixNano())
	return e.limiter
}

func (k *RateLimiter) cleanup() {
	ticker := time.NewTicker(k.cleanUpIn)
	defer ticker.Stop()
	
	for {
		select {
		case <-k.stopChan:
			return
		case <-ticker.C:
			now := time.Now().UnixNano()
			k.limiters.Range(func(key, value any) bool {
				e := value.(*limiterEntry)
				if now-e.lastSeen.Load() > int64(k.cleanUpIn) {
					k.limiters.Delete(key)
				}
				return true
			})
		}
	}
}

func (k *RateLimiter) Stop() {
	close(k.stopChan)
}

func GetKey(ctx context.Context) (string, error) {
	key := api.MetaFp(ctx)
	if key == "" {
		key = api.MetaIp(ctx)
	}
	return key, nil
}

// UnaryServerInterceptor 根据 IP 限流
func (k *RateLimiter) UnaryServerInterceptor(getKeyFunc ...func(ctx context.Context) (string, error)) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		f := GetKey
		if len(getKeyFunc) > 0 {
			f = getKeyFunc[0]
		}
		key, err := f(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get rate limit key: %v", err)
		}

		limiter := k.getLimiter(key)
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded for %s", key)
		}

		return handler(ctx, req)
	}
}
