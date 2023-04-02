package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
)

type Redis struct {
	rdc *redis.Client
}

func New(app *stargo.App) *Redis {

	c := &Redis{
		rdc: app.GetRedis().GetInstance(),
	}

	// go func() {
	// 	t := time.NewTicker(time.Second * 5) //TODO，传入配置，interval
	// 	defer t.Stop()
	// 	for {
	// 		<-t.C
	// 		c.clear(context.Background())
	// 	}
	// }()

	return c
}

func (e *Redis) Get(ctx context.Context, key string) (any, error) {
	_, err := e.rdc.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Put stores a key-value pair into cache.
func (e *Redis) Put(ctx context.Context, key string, value any) error {
	return nil
}

// Delete removes a key from cache.
func (e *Redis) Delete(ctx context.Context, key string) error {
	return nil
}
func (e *Redis) Scan(ctx context.Context, key string, data any) error {
	if err := e.rdc.Get(ctx, key).Scan(data); err != nil {
		return err
	}
	return nil
}

// 过期清除，好像用不着
func (e *Redis) clear(ctx context.Context) {
	iter := e.rdc.Scan(ctx, 0, "", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		d, err := e.rdc.TTL(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		if d == -1 { // -1 means no TTL
			if err := e.rdc.Del(ctx, key).Err(); err != nil {
				panic(err)
			}
		}
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}

}
