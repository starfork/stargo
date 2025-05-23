package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	cache "github.com/starfork/gorm-cache"
)

type redisCacher struct {
	rdc *redis.Client
}

func NewRedisCacher(rdc *redis.Client) *redisCacher {
	return &redisCacher{rdc: rdc}
}

func (c *redisCacher) Get(ctx context.Context, key string, q *cache.Query[any]) (*cache.Query[any], error) {
	res, err := c.rdc.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := q.Unmarshal([]byte(res)); err != nil {
		return nil, err
	}

	return q, nil
}

func (c *redisCacher) Store(ctx context.Context, key string, val *cache.Query[any]) error {
	res, err := val.Marshal()
	if err != nil {
		return err
	}

	c.rdc.Set(ctx, key, res, 300*time.Second) // Set proper cache time
	return nil
}

func (c *redisCacher) Invalidate(ctx context.Context) error {
	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = c.rdc.Scan(ctx, cursor, fmt.Sprintf("%s*", cache.IdentifierPrefix), 0).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if _, err := c.rdc.Del(ctx, keys...).Result(); err != nil {
			return err
		}
	}
	return nil
}
