package etcd

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/util/ustring"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

// const KeyPrefix = "stargo_registry"
const Scheme = "etcd"

//const Org = "stargo"

type Registry struct {
	name string
	org  string
	cli  *clientv3.Client

	em   endpoints.Manager
	conf *naming.Config

	ctx    context.Context
	cancel context.CancelFunc
	leaseID clientv3.LeaseID
	mu      sync.Mutex
}

func init() {
	naming.RegisterRegistry(Scheme, func(conf *naming.Config) (naming.Registry, error) {
		return NewRegistry(conf)
	})
}

func NewRegistry(conf *naming.Config) (naming.Registry, error) {

	cli, err := newClient(conf)
	if err != nil {
		return nil, err
	}
	org := ustring.OrString("stargo", conf.Org)
	em, err := endpoints.NewManager(cli, org)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Registry{
		org:    org,
		cli:    cli,
		name:   Scheme,
		em:     em,
		conf:   conf,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}
func (e *Registry) key(svc naming.Service) string {
	key := e.org + "/" + svc.Name + "/" + svc.Addr
	return key
}

// serviceMetadata stores Weight and Version as etcd endpoint metadata
type serviceMetadata struct {
	Version string `json:"ver,omitempty"`
	Weight  int64  `json:"wt,omitempty"`
}

func (e *Registry) Register(svc naming.Service) error {
	meta := serviceMetadata{
		Version: svc.Version,
		Weight:  svc.Weight,
	}
	if svc.Weight <= 0 {
		meta.Weight = 100
	}
	metaBytes, _ := json.Marshal(meta)
	p := endpoints.Endpoint{
		Addr:     svc.Addr,
		Metadata: json.RawMessage(metaBytes),
	}
	opts := []clientv3.OpOption{}
	if e.conf.Ttl > 0 {
		lease, err := e.cli.Grant(e.cli.Ctx(), e.conf.Ttl)
		if err != nil {
			logger.DefaultLogger.Errorf("etcd lease grant failed: %v", err)
			return err
		}
		opts = append(opts, clientv3.WithLease(lease.ID))
		e.mu.Lock()
		e.leaseID = lease.ID
		e.mu.Unlock()

		go e.keepAlive(lease.ID, svc)
	}
	return e.em.AddEndpoint(e.cli.Ctx(), e.key(svc), p, opts...)
}

func (e *Registry) Deregister(svc naming.Service) error {
	e.mu.Lock()
	if e.leaseID != 0 {
		_, err := e.cli.Revoke(e.ctx, e.leaseID)
		if err != nil {
			logger.DefaultLogger.Errorf("etcd lease revoke failed: %v", err)
		}
		e.leaseID = 0
	}
	e.mu.Unlock()
	return e.em.DeleteEndpoint(e.cli.Ctx(), e.key(svc))
}

func (e *Registry) List(name string) []naming.Service {
	prefix := e.org + "/" + name + "/"
	resp, err := e.cli.Get(e.cli.Ctx(), prefix, clientv3.WithPrefix())
	if err != nil {
		logger.DefaultLogger.Errorf("List failed: %v", err)
		return nil
	}
	var services []naming.Service
	for _, kv := range resp.Kvs {
		parts := strings.Split(string(kv.Key), "/")
		if len(parts) < 3 {
			continue
		}
		svc := naming.Service{
			Org:  parts[0],
			Name: parts[1],
			Addr: parts[2],
		}
		if len(kv.Value) > 0 {
			var meta serviceMetadata
			if err := json.Unmarshal(kv.Value, &meta); err == nil {
				svc.Version = meta.Version
				svc.Weight = meta.Weight
			}
		}
		services = append(services, svc)
	}
	return services
}

func (e *Registry) Scheme() string {
	return e.name
}

func (e *Registry) keepAlive(leaseID clientv3.LeaseID, svc naming.Service) {
	ch, err := e.cli.KeepAlive(e.ctx, leaseID)
	if err != nil {
		logger.DefaultLogger.Errorf("etcd keepalive failed: %v", err)
		return
	}

	for {
		select {
		case <-e.ctx.Done():
			return
		case resp, ok := <-ch:
			if !ok {
				logger.DefaultLogger.Warnf("etcd keepalive channel closed, attempting re-register")
				e.reRegister(svc)
				return
			}
			if resp == nil {
				logger.DefaultLogger.Warnf("etcd keepalive response is nil")
			}
		}
	}
}

func (e *Registry) reRegister(svc naming.Service) {
	p := endpoints.Endpoint{
		Addr: svc.Addr,
	}
	
	// Exponential backoff: 1s, 2s, 4s, 8s, 16s, 30s (max)
	backoff := time.Second
	maxBackoff := 30 * time.Second
	
	for {
		select {
		case <-e.ctx.Done():
			return
		default:
		}
		
		if e.conf.Ttl > 0 {
			lease, err := e.cli.Grant(e.cli.Ctx(), e.conf.Ttl)
			if err != nil {
				logger.DefaultLogger.Errorf("etcd re-register grant failed: %v, retrying in %v", err, backoff)
				time.Sleep(backoff)
				backoff = min(backoff*2, maxBackoff)
				continue
			}
			
			// Create fresh opts each iteration to avoid accumulating leases
			opts := []clientv3.OpOption{clientv3.WithLease(lease.ID)}
			e.mu.Lock()
			e.leaseID = lease.ID
			e.mu.Unlock()
			
			if err := e.em.AddEndpoint(e.cli.Ctx(), e.key(svc), p, opts...); err != nil {
				logger.DefaultLogger.Errorf("etcd re-register failed: %v, retrying in %v", err, backoff)
				time.Sleep(backoff)
				backoff = min(backoff*2, maxBackoff)
				continue
			}
			
			logger.DefaultLogger.Infof("etcd re-register successful")
			// Restart keepalive
			go e.keepAlive(lease.ID, svc)
			return
		}
	}
}

func (e *Registry) Close() error {
	e.cancel()
	e.mu.Lock()
	if e.leaseID != 0 {
		_, err := e.cli.Revoke(e.ctx, e.leaseID)
		if err != nil {
			logger.DefaultLogger.Errorf("etcd lease revoke failed: %v", err)
		}
	}
	e.mu.Unlock()
	return e.cli.Close()
}
