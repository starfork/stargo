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
	rdc *redis.Client
}

func New(s store.Store) cache.Cache {

	c := &Redis{
		rdc: s.Instance().(*redis.Client),
	}

	return c
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
	return nil, nil
}

func (e *Redis) Put(ctx context.Context, key string, value any, timeout ...time.Duration) error {
	var expr time.Duration = -1
	if len(timeout) > 0 {
		expr = timeout[0]
	}
	rs := e.rdc.Set(ctx, key, value, expr)
	return rs.Err()
}

// Delete removes a key from cache.
func (e *Redis) Delete(ctx context.Context, key string) error {
	e.rdc.Del(ctx, key).Result()
	return nil
}
func (e *Redis) Clear(ctx context.Context, key string) error {
	iter := e.rdc.Scan(ctx, 0, key+"*", 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("del key ", iter.Val())
		_, err := e.rdc.Del(ctx, iter.Val()).Result()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func (e *Redis) IsExist(ctx context.Context, key string) (bool, error) {
	return false, nil
}
func (e *Redis) ClearAll(ctx context.Context) error {
	return nil
}

func (e *Redis) Incr(ctx context.Context, key string) error {
	return nil
}

// Decrement a cached int value by key, as a counter.
func (e *Redis) Decr(ctx context.Context, key string) error {
	return nil
}

func (e *Redis) Expire(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// func (e *Redis) Scan(ctx context.Context, key string, data any) error {
// 	if err := e.rdc.Get(ctx, key).Scan(data); err != nil {
// 		return err
// 	}
// 	return nil
// }
