package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/service"
	ssredis "github.com/starfork/stargo/store/redis"
	"google.golang.org/grpc/resolver"
)

const Scheme = "redis"

type Resolver struct {
	rdc  *redis.Client
	name string
	ctx  context.Context
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
	}

	resolver.Register(r)
	return r
}

func (e *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var address []resolver.Address
	key := e.name + "_" + target.URL.Host
	fmt.Println(key)
	rs := e.rdc.SMembers(e.ctx, key)
	if rs.Err() != nil {
		fmt.Println(rs.Err())
		return nil, rs.Err()
	}

	fmt.Printf("%+v", rs.Val())

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
	key := e.name + "_" + name
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
