package config

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigManager interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Watch(key string, onChange func(value string))
}

type DefaultConfigManager struct {
	client     *clientv3.Client
	config     map[string]string
	watchChans map[string]clientv3.WatchChan
}

func NewDefaultConfigManager(endpoints []string) (ConfigManager, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
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

func (cm *DefaultConfigManager) Watch(key string, onChange func(value string)) {
	watchChan := cm.client.Watch(context.Background(), key)
	cm.watchChans[key] = watchChan
	go func() {
		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				if ev.Type == clientv3.EventTypePut {
					value := string(ev.Kv.Value)
					cm.config[key] = value
					onChange(value)
				}
			}
		}
	}()
}
