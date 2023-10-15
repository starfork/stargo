package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/service"
	ssredis "github.com/starfork/stargo/store/redis"
)

type Registry struct {
	rdc  *redis.Client
	name string
	ctx  context.Context
	conf *config.Registry
}

func NewRegistry(conf *config.Registry) *Registry {

	r := ssredis.Connect(&config.Config{
		Redis: &config.RedisConfig{
			Addr: conf.Addr,
			Auth: conf.Auth,
			//Num:  conf.Num,
		},
	})
	return &Registry{
		rdc: r.GetInstance(),
		//rdc:  app.GetRedis().GetInstance(),
		ctx:  context.Background(),
		name: Scheme,
		conf: conf,
	}
}

// /"stargo_registry_[or]_[service]"
func (e *Registry) key(name string) string {
	return KeyPrefix + "_" + e.conf.Org + "_" + name
}
func (e *Registry) Register(svc service.Service) error {
	key := e.key(svc.Name)
	err := e.rdc.SAdd(context.TODO(), key, svc.Addr).Err()
	return err
}

func (e *Registry) UnRegister(svc service.Service) error {

	key := e.key(svc.Name)
	err := e.rdc.SRem(e.ctx, key, svc.Addr).Err()
	return err
}

func (e *Registry) List(name string) []service.Service {
	key := e.key(name)

	rs := e.rdc.SMembers(e.ctx, key)
	data := []service.Service{}
	for _, v := range rs.Val() {
		data = append(data, service.Service{
			Name: name,
			Addr: v,
		})
	}

	return data
}

func (e *Registry) Scheme() string {
	return e.name
}
