package etcd

import (
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

	em   endpoints.Manager
	conf *config.Registry
	//services []service.Service
}

func NewRegistry(conf *config.Registry) *Registry {
	cli := newClient(conf)

	//defer cli.Close()

	em, err := endpoints.NewManager(cli, conf.Org)
	if err != nil {
		panic(err)
	}

	return &Registry{
		cli:  cli,
		name: Scheme,
		em:   em,
		conf: conf,
	}
}
func (e *Registry) key(name string) string {
	return e.conf.Org + "/" + name
}

func (e *Registry) Register(svc service.Service) error {

	p := endpoints.Endpoint{
		Addr: svc.Addr,
	}
	key := e.key(svc.Name)
	return e.em.AddEndpoint(e.cli.Ctx(), key, p)
}

func (e *Registry) UnRegister(svc service.Service) error {
	key := e.key(svc.Name)
	return e.em.DeleteEndpoint(e.cli.Ctx(), key)
}

func (e *Registry) List(name string) []service.Service {

	return nil
}

func (e *Registry) Scheme() string {
	return e.name
}
