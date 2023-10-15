package naming

import (
	"strings"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming/etcd"
	"github.com/starfork/stargo/naming/redis"
	"github.com/starfork/stargo/service"
)

type Registry interface {
	Scheme() string

	Register(service service.Service) error
	UnRegister(service service.Service) error
	//返回服务
	List(name string) []service.Service
}

func NewRegistry(conf *config.Registry) Registry {
	name := strings.ToLower(conf.Name)
	if name == "redis" {
		return redis.NewRegistry(conf)
	}
	if name == "etcd" {
		return etcd.NewRegistry(conf)
	}
	panic("unknow registry")
}
