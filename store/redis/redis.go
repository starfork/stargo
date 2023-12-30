package redis

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/util/ustring"
)

type Redis struct {
	rdc *redis.Client
	c   *config.StoreConfig
}

func NewRedis(config *config.StoreConfig) store.Store {

	return &Redis{
		c: config,
	}
}

func (e *Redis) Connect(conf ...*config.Config) {
	c := e.c
	if len(conf) > 0 && conf[0].Redis != nil {
		c = conf[0].Redis
	}
	c.Host = ustring.Or(c.Host, os.Getenv("REDIS_HOST"))
	c.Auth = ustring.Or(c.Auth, os.Getenv("REDIS_AUTH"))
	c.Num = ustring.Int(ustring.OrString(strconv.Itoa(c.Num), os.Getenv("REDIS_NUM")))

	rdc := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		DB:       c.Num,
		Password: c.Auth,
	})

	if _, err := rdc.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	e.rdc = rdc
}

func (e *Redis) GetInstance(conf ...*config.Config) *redis.Client {
	if len(conf) > 0 {
		e.Connect()
		return e.rdc
	}
	return e.rdc
}

func (e *Redis) Close() {
	if e.rdc != nil {
		e.rdc.Close()
	}
}
