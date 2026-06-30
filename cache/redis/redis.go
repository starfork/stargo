package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/cache"
	"github.com/starfork/stargo/store"
)

type Redis struct {
	rdc  *redis.Client
	name string
}

func New(s store.Store, name string) (cache.Cache, error) {
	instance, err := s.InstanceE()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis instance: %w", err)
	}
	
	rdc, ok := instance.(*redis.Client)
	if !ok || rdc == nil {
		return nil, fmt.Errorf("invalid redis instance type")
	}

	c := &Redis{
		rdc:  rdc,
		name: name,
	}

	return c, nil
}

func (e *Redis) Get(ctx context.Context, key string) (any, error) {

	rs, err := e.rdc.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// Batch of get
func (e *Redis) Fetch(ctx context.Context, key []string) ([]any, error) {
	if len(key) == 0 {
		return nil, nil
	}
	
	rs, err := e.rdc.MGet(ctx, key...).Result()
	if err != nil {
		return nil, err
	}
	
	result := make([]any, len(rs))
	for i, v := range rs {
		result[i] = v
	}
	return result, nil
}

func (e *Redis) Put(ctx context.Context, key string, value any, timeout ...time.Duration) error {
	var expr time.Duration = 0 // 0 means no expiration
	if len(timeout) > 0 {
		expr = timeout[0]
	}
	rs := e.rdc.Set(ctx, key, value, expr)
	return rs.Err()
}

// Delete removes a key from cache.
func (e *Redis) Delete(ctx context.Context, key string) error {
	return e.rdc.Del(ctx, key).Err()
}
func (e *Redis) Clear(ctx context.Context, key string) error {
	iter := e.rdc.Scan(ctx, 0, key+"*", 0).Iterator()
	for iter.Next(ctx) {
		_, err := e.rdc.Del(ctx, iter.Val()).Result()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func (e *Redis) IsExist(ctx context.Context, key string) (bool, error) {
	rs, err := e.rdc.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return rs > 0, nil
}
func (e *Redis) ClearAll(ctx context.Context) error {
	// Use service prefix to avoid deleting keys from other services
	pattern := e.name + "*"
	if e.name == "" {
		pattern = "*"
	}
	iter := e.rdc.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := e.rdc.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func (e *Redis) Incr(ctx context.Context, key string) error {
	return e.rdc.Incr(ctx, key).Err()
}

// Decrement a cached int value by key, as a counter.
func (e *Redis) Decr(ctx context.Context, key string) error {
	return e.rdc.Decr(ctx, key).Err()
}

func (e *Redis) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		// If ttl is 0 or negative, remove the key
		result, err := e.rdc.Del(ctx, key).Result()
		return result > 0, err
	}
	return e.rdc.Expire(ctx, key, ttl).Result()
}

// func (e *Redis) Scan(ctx context.Context, key string, data any) error {
// 	if err := e.rdc.Get(ctx, key).Scan(data); err != nil {
// 		return err
// 	}
// 	return nil
// }
