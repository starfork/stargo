package redis

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/util/ustring"
)

type Redis struct {
	rdc *redis.Client
}

func Connect(config *config.Config) *Redis {
	c := config.Redis
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
	return &Redis{
		rdc: rdc,
	}
}

func (e *Redis) GetInstance(conf ...*config.Config) *redis.Client {
	if len(conf) > 0 {
		rs := Connect(conf[0])
		return rs.rdc
	}
	return e.rdc
}

func (e *Redis) Close() {
	if e.rdc != nil {
		e.rdc.Close()
	}
}
