package redis

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/util/ustring"
)

type Redis struct {
	rdc *redis.Client
	c   *store.Config
}

func NewRedis(config *store.Config) store.Store {
	return &Redis{
		c: config,
	}
}

func (e *Redis) Connect(conf ...*store.Config) {

	c := e.c
	if len(conf) > 0 {
		c = conf[0]
	}
	c.Host = ustring.Or(c.Host, os.Getenv("REDIS_HOST"))
	c.Auth = ustring.Or(c.Auth, os.Getenv("REDIS_AUTH"))
	c.Num = ustring.Int(ustring.OrString(strconv.Itoa(c.Num), os.Getenv("REDIS_NUM")))

	rdc := redis.NewClient(&redis.Options{
		Addr:     c.Host + ":" + c.Port,
		DB:       c.Num,
		Password: c.Auth,
	})

	if _, err := rdc.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	e.rdc = rdc
}

func (e *Redis) Instance(conf ...*store.Config) any {
	if len(conf) > 0 {
		e.Connect(conf...)
		return e.rdc
	}
	if e.rdc == nil {
		e.Connect()
	}
	return e.rdc
}

// 集群client
func (e *Redis) GetCluster(conf ...*store.Config) *redis.ClusterClient {
	return nil
}

func (e *Redis) Close() {
	if e.rdc != nil {
		e.rdc.Close()
	}
}
