package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/cache"
)

type Redis struct {
	rdc *redis.Client
}

func New(app *stargo.App) *Redis {

	return &Redis{
		rdc: app.GetRedis().GetInstance(),
	}
}

func (e *Redis) Get(ctx context.Context, key string) (*cache.Item, error) {
	_, err := e.rdc.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Put stores a key-value pair into cache.
func (e *Redis) Put(ctx context.Context, key string, value *cache.Item) error {
	return nil
}

// Delete removes a key from cache.
func (e *Redis) Delete(ctx context.Context, key string) error {
	return nil
}
