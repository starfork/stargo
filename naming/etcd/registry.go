package etcd

import (
	"context"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/service"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

const KeyPrefix = "stargo_registry"
const Scheme = "etcd"

type Registry struct {
	name string
	cli  *clientv3.Client

	em       endpoints.Manager
	ctx      context.Context
	conf     *config.Registry
	services []service.Service
}

func NewRegistry(conf *config.Registry) *Registry {
	cli, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		panic(err)
	}
	return &Registry{
		cli:  cli,
		ctx:  context.Background(),
		name: Scheme,
		conf: conf,
	}
}

func (e *Registry) Register(svc service.Service) error {

	return nil
}

func (e *Registry) UnRegister(svc service.Service) error {

	return nil
}

func (e *Registry) List(name string) []service.Service {

	return nil
}

func (e *Registry) Scheme() string {
	return e.name
}
