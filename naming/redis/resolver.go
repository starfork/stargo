package redis

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/service"
	sredis "github.com/starfork/stargo/store/redis"
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
	rds := sredis.NewRedis(&config.StoreConfig{
		Host: conf.Host,
		Auth: conf.Auth,
	}).(*sredis.Redis)

	r := &Resolver{
		name: Scheme,
		rdc:  rds.GetInstance(),
		ctx:  context.Background(),
		conf: conf,
	}

	resolver.Register(r)
	return r
}

// /"stargo_registry_[or]_[service]"
func (e *Resolver) key(name string) string {
	return Scheme + "_" + e.conf.Org + "_" + name
}
func (e *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var address []resolver.Address
	key := e.key(target.URL.Host)

	rs := e.rdc.SMembers(e.ctx, key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	if len(rs.Val()) == 0 {
		return nil, errors.New(" 无可用注册服务: " + target.URL.Host)
	}

	for _, v := range rs.Val() {
		t := strings.Split(v, ":")
		if e.conf.Environment == "debug" {
			v = ":" + t[1] //测试环境，只要端口号
		}
		if e.conf.Environment == "docker" {
			v = "host.docker.internal:" + t[1]
		}
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
