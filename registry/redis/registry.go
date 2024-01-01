package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/registry"
	ssredis "github.com/starfork/stargo/store/redis"
)

const KeyPrefix = "stargo_registry"
const Scheme = "redis"

type Registry struct {
	rdc  *redis.Client
	name string
	ctx  context.Context
	conf *config.Registry
}

func NewRegistry(conf *config.Registry) *Registry {
	rds := ssredis.NewRedis(&config.StoreConfig{
		Host: conf.Host,
		Auth: conf.Auth,
	}).(*ssredis.Redis)
	return &Registry{
		rdc: rds.GetInstance(),
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
func (e *Registry) Register(svc registry.Service) error {
	key := e.key(svc.Name)
	err := e.rdc.SAdd(context.TODO(), key, svc.Addr).Err()
	return err
}

func (e *Registry) UnRegister(svc registry.Service) error {

	key := e.key(svc.Name)
	err := e.rdc.SRem(e.ctx, key, svc.Addr).Err()
	return err
}

func (e *Registry) List(name string) []registry.Service {
	key := e.key(name)

	rs := e.rdc.SMembers(e.ctx, key)
	data := []registry.Service{}
	for _, v := range rs.Val() {
		data = append(data, registry.Service{
			Name: name,
			Addr: v,
		})
	}

	return data
}

func (e *Registry) Scheme() string {
	return e.name
}
