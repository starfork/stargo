package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/cache"
	sredis "github.com/starfork/stargo/store/redis"
)

type Redis struct {
	rdc *redis.Client
}

func New(app *stargo.App) cache.Cache {

	c := &Redis{
		rdc: app.Store("redis").(*sredis.Redis).GetInstance(),
	}

	return c
}

func (e *Redis) Get(ctx context.Context, key string) (any, error) {
	_, err := e.rdc.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Batch of get
func (e *Redis) Fetch(ctx context.Context, key []string) ([]any, error) {
	return nil, nil
}

func (e *Redis) Put(ctx context.Context, key string, value any, timeout ...time.Duration) error {
	return nil
}

// Delete removes a key from cache.
func (e *Redis) Delete(ctx context.Context, key string) error {
	return nil
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

// func (e *Redis) Scan(ctx context.Context, key string, data any) error {
// 	if err := e.rdc.Get(ctx, key).Scan(data); err != nil {
// 		return err
// 	}
// 	return nil
// }

// 过期清除，好像用不着
// func (e *Redis) Clear(ctx context.Context) {
// 	iter := e.rdc.Scan(ctx, 0, "", 0).Iterator()

// 	for iter.Next(ctx) {
// 		key := iter.Val()

// 		d, err := e.rdc.TTL(ctx, key).Result()
// 		if err != nil {
// 			panic(err)
// 		}
// 		if d == -1 { // -1 means no TTL
// 			if err := e.rdc.Del(ctx, key).Err(); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	if err := iter.Err(); err != nil {
// 		panic(err)
// 	}

// }
