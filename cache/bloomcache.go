package cache

import (
	"context"
	"time"

	"github.com/starfork/stargo/cache/bloom"
)

type bloomCache struct {
	Cache
	bf *bloom.BloomFilter
}

func NewBloomCache(inner Cache, bf *bloom.BloomFilter) Cache {
	return &bloomCache{Cache: inner, bf: bf}
}

func (b *bloomCache) Get(ctx context.Context, key string) (interface{}, error) {
	if !b.bf.MightContain(key) {
		return nil, nil
	}
	return b.Cache.Get(ctx, key)
}

func (b *bloomCache) Put(ctx context.Context, key string, value interface{}, timeout ...time.Duration) error {
	b.bf.Add(key)
	return b.Cache.Put(ctx, key, value, timeout...)
}
