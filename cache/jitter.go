package cache

import (
	"context"
	"math/rand"
	"time"
)

type jitterCache struct {
	Cache
	jitterPercent float64
}

func NewJitterCache(inner Cache, jitterPercent float64) Cache {
	if jitterPercent <= 0 {
		jitterPercent = 0.2
	}
	return &jitterCache{Cache: inner, jitterPercent: jitterPercent}
}

func (j *jitterCache) Put(ctx context.Context, key string, value interface{}, timeout ...time.Duration) error {
	if len(timeout) > 0 && timeout[0] > 0 {
		ttl := jitterTTL(timeout[0], j.jitterPercent)
		return j.Cache.Put(ctx, key, value, ttl)
	}
	return j.Cache.Put(ctx, key, value, timeout...)
}

func jitterTTL(base time.Duration, jitterPercent float64) time.Duration {
	delta := time.Duration(float64(base) * jitterPercent)
	minTTL := base - delta
	maxTTL := base + delta
	if minTTL <= 0 {
		minTTL = base / 2
	}
	return minTTL + time.Duration(rand.Int63n(int64(maxTTL-minTTL+1)))
}
