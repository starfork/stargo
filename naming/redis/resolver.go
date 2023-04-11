package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/service"
	ssredis "github.com/starfork/stargo/store/redis"
	"google.golang.org/grpc/resolver"
)

const KeyPrefix = "stargo_registry"
const Scheme = "redis"

type Resolver struct {
	rdc  *redis.Client
	name string
	ctx  context.Context
	conf *config.Registry
}

func NewResolver(conf *config.Registry) resolver.Builder {
	rds := ssredis.Connect(&config.ServerConfig{
		Redis: &config.RedisConfig{
			Addr: conf.Addr,
			//Auth: conf.Auth,
			//Num:  conf.Num,
		},
	})

	r := &Resolver{
		name: Scheme,
		rdc:  rds.GetInstance(),
		ctx:  context.Background(),
		conf: conf,
	}

	resolver.Register(r)
	return r
}

// stargo_registryredis[xxx]abc
func (e *Resolver) key(name string) string {
	return KeyPrefix + "_" + e.conf.Org + "_" + name
}
func (e *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var address []resolver.Address
	key := e.key(target.URL.Host)
	rs := e.rdc.SMembers(e.ctx, key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}

	for _, v := range rs.Val() {
		address = append(address, resolver.Address{
			Addr: v,
		})
	}
	cc.UpdateState(resolver.State{Addresses: address})
	//resolver
	return &nopResolver{}, nil
}

func (e *Resolver) List(name string) []service.Service {

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

func (b *Resolver) Scheme() string {
	return Scheme
}

type nopResolver struct {
}

func (*nopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*nopResolver) Close() {}
