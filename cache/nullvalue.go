package cache

import (
	"context"
	"time"
)

const nullValue = "__stargo_cache_null__"

type nullValueCache struct {
	Cache
	nullTTL time.Duration
}

func NewNullValueCache(inner Cache, nullTTL time.Duration) Cache {
	if nullTTL <= 0 {
		nullTTL = 30 * time.Second
	}
	return &nullValueCache{Cache: inner, nullTTL: nullTTL}
}

func (n *nullValueCache) GetOrLoad(ctx context.Context, key string, dataTTL time.Duration, loader func() (interface{}, bool)) (interface{}, error) {
	val, err := n.Cache.Get(ctx, key)
	if err == nil {
		if val == nullValue {
			return nil, nil
		}
		return val, nil
	}

	data, found := loader()
	if !found {
		n.Cache.Put(ctx, key, nullValue, n.nullTTL)
		return nil, nil
	}
	n.Cache.Put(ctx, key, data, dataTTL)
	return data, nil
}
