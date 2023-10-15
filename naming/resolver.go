package naming

import (
	"strings"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming/etcd"
	"github.com/starfork/stargo/naming/redis"
	"google.golang.org/grpc/resolver"
)

type Resolver interface {
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)
	Scheme() string
}

func NewResolver(conf *config.Registry) Resolver {
	name := strings.ToLower(conf.Name)
	if name == "redis" {
		return redis.NewResolver(conf)
	}
	if name == "etcd" {
		return etcd.NewResolver(conf)
	}
	panic("unknow resolver")
}
