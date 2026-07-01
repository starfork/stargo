package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type singleflightCache struct {
	Cache
	sg singleflight.Group
}

func NewSingleflightCache(inner Cache) Cache {
	return &singleflightCache{Cache: inner, sg: singleflight.Group{}}
}

func (s *singleflightCache) GetOrLoad(ctx context.Context, key string, ttl time.Duration, loader func() (interface{}, error)) (interface{}, error) {
	if val, err := s.Cache.Get(ctx, key); err == nil && val != nil {
		return val, nil
	}

	v, err, _ := s.sg.Do(key, func() (interface{}, error) {
		val, err := loader()
		if err != nil {
			return nil, err
		}
		if putErr := s.Cache.Put(ctx, key, val, ttl); putErr != nil {
			return nil, putErr
		}
		return val, nil
	})
	return v, err
}
