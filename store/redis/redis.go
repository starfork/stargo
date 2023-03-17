package redis

import (
	"public/pkg/config"

	"github.com/go-redis/redis"
)

var (
	//Rdc rediscli
	rdc *redis.Client
)

type Redis struct {
	client *redis.Client
}

func Connect(config *config.ServerConfig) *Redis {
	c := config.Redis
	rdc = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		DB:       c.Num,
		Password: c.Auth,
	})

	_, err := rdc.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &Redis{
		client: rdc,
	}
}

func (e *Redis) GetInstance() *redis.Client {
	return e.client
}

func (e *Redis) Close() {
	if e.client != nil {
		e.client.Close()
	}
}
