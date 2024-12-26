package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigManager interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Watch(key string, onChange func(key, value string))
}

type DefaultConfigManager struct {
	client     *clientv3.Client
	config     map[string]string
	watchChans map[string]clientv3.WatchChan
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
		watchChans: make(map[string]clientv3.WatchChan),
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
	watchChan := cm.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	go func() {
		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				if ev.Type == clientv3.EventTypePut {
					fullKey := string(ev.Kv.Key)
					value := string(ev.Kv.Value)
					subKey := strings.TrimPrefix(fullKey, prefix)
					subKey = strings.TrimLeft(subKey, "/")
					onChange(subKey, value)
				}
			}
		}
	}()
}
