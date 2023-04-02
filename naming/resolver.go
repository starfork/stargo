package naming

import (
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming/redis"
	"google.golang.org/grpc/resolver"
)

type Resolver interface {
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)
	Scheme() string
}

func NewResolver(conf *config.Registry) Resolver {
	if conf.Name == "redis" {
		return redis.NewResolver(conf)
	}
	return nil
}
