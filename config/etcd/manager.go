package etcdconfig

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/starfork/stargo/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigManager interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Watch(key string, onChange func(key, value string))
	Close() error
}

type watchEntry struct {
	cancel    context.CancelFunc
	onChange  func(k, v string)
}

type DefaultConfigManager struct {
	client     *clientv3.Client
	config     map[string]string
	watchChans map[string]*watchEntry
	mu         sync.RWMutex
}

func NewDefaultConfigManager(endpoints string) (ConfigManager, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(endpoints, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &DefaultConfigManager{
		client:     cli,
		config:     make(map[string]string),
		watchChans: make(map[string]*watchEntry),
	}, nil
}

func (cm *DefaultConfigManager) Get(key string) (string, error) {
	resp, err := cm.client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

func (cm *DefaultConfigManager) Set(key, value string) error {
	_, err := cm.client.Put(context.Background(), key, value)
	return err
}

func (cm *DefaultConfigManager) Watch(prefix string, onChange func(k, v string)) {
	cm.mu.Lock()
	
	// Cancel existing watch if any
	if entry, ok := cm.watchChans[prefix]; ok {
		entry.cancel()
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	entry := &watchEntry{
		cancel:   cancel,
		onChange: onChange,
	}
	cm.watchChans[prefix] = entry
	cm.mu.Unlock()
	
	go cm.watchLoop(ctx, prefix, onChange)
}

func (cm *DefaultConfigManager) watchLoop(ctx context.Context, prefix string, onChange func(k, v string)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		
		watchChan := cm.client.Watch(ctx, prefix, clientv3.WithPrefix())
		cm.processWatchChan(ctx, prefix, watchChan, onChange)
		
		// Wait before reconnecting
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			logger.DefaultLogger.Infof("reconnecting watch for prefix: %s", prefix)
		}
	}
}

func (cm *DefaultConfigManager) processWatchChan(ctx context.Context, prefix string, watchChan clientv3.WatchChan, onChange func(k, v string)) {
	for watchResp := range watchChan {
		if watchResp.Err() != nil {
			logger.DefaultLogger.Errorf("watch error for prefix %s: %v", prefix, watchResp.Err())
			return
		}
		
		for _, ev := range watchResp.Events {
			fullKey := string(ev.Kv.Key)
			subKey := strings.TrimPrefix(fullKey, prefix)
			subKey = strings.TrimLeft(subKey, "/")
			
			switch ev.Type {
			case clientv3.EventTypePut:
				value := string(ev.Kv.Value)
				onChange(subKey, value)
			case clientv3.EventTypeDelete:
				onChange(subKey, "")
			}
		}
	}
}

func (cm *DefaultConfigManager) Close() error {
	cm.mu.Lock()
	for prefix, entry := range cm.watchChans {
		entry.cancel()
		delete(cm.watchChans, prefix)
	}
	cm.mu.Unlock()
	
	if cm.client != nil {
		return cm.client.Close()
	}
	return nil
}
