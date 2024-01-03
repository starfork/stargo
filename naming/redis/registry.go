package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
	ssredis "github.com/starfork/stargo/store/redis"
)

type Registry struct {
	rdc  *redis.Client
	name string
	ctx  context.Context
	conf *naming.Config
}

func NewRegistry(conf *naming.Config) *Registry {
	rds := ssredis.NewRedis(&store.Config{
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
func (e *Registry) Register(svc naming.Service) error {
	key := e.key(svc.Name)
	err := e.rdc.SAdd(context.TODO(), key, svc.Addr).Err()
	return err
}

func (e *Registry) UnRegister(svc naming.Service) error {

	key := e.key(svc.Name)
	err := e.rdc.SRem(e.ctx, key, svc.Addr).Err()
	return err
}

func (e *Registry) List(name string) []naming.Service {
	key := e.key(name)

	rs := e.rdc.SMembers(e.ctx, key)
	data := []naming.Service{}
	for _, v := range rs.Val() {
		data = append(data, naming.Service{
			Name: name,
			Addr: v,
		})
	}

	return data
}

func (e *Registry) Scheme() string {
	return e.name
}
