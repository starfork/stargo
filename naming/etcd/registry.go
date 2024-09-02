package etcd

import (
	"github.com/starfork/stargo/naming"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

const KeyPrefix = "stargo_registry"
const Scheme = "etcd"

type Registry struct {
	name string
	cli  *clientv3.Client

	em   endpoints.Manager
	conf *naming.Config
	//services []registry.Service
}

func NewRegistry(conf *naming.Config) naming.Registry {
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

func (e *Registry) Register(svc naming.Service) error {

	p := endpoints.Endpoint{
		Addr: svc.Addr,
	}
	key := e.key(svc.Name)
	return e.em.AddEndpoint(e.cli.Ctx(), key, p)
}

func (e *Registry) Deregister(svc naming.Service) error {
	key := e.key(svc.Name)
	return e.em.DeleteEndpoint(e.cli.Ctx(), key)
}

func (e *Registry) List(name string) []naming.Service {

	return nil
}

func (e *Registry) Scheme() string {
	return e.name
}
