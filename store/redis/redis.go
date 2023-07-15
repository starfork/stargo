package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
)

type Redis struct {
	rdc *redis.Client
}

func Connect(config *config.ServerConfig) *Redis {
	c := config.Redis
	rdc := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		DB:       c.Num,
		Password: c.Auth,
	})

	if _, err := rdc.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &Redis{
		rdc: rdc,
	}
}

func (e *Redis) GetInstance() *redis.Client {
	return e.rdc
}

func (e *Redis) Close() {
	if e.rdc != nil {
		e.rdc.Close()
	}
}
