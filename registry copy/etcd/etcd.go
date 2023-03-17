package etcd

import (
	"public/pkg/registry"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type etcdRegistry struct {
	service registry.Service
	em      endpoints.Manager
	cli     *clientv3.Client
}

// NewRegister create a register
func NewEtcdRegistry(addrs []string, ttl int) (registry.Registry, error) {

	//var err error
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: time.Second * time.Duration(ttl),
	})
	if err != nil {
		return nil, err
	}

	r := &etcdRegistry{
		cli: cli,
	}

	return r, nil
}

func (e *etcdRegistry) Register(service registry.Service) error {

	name := e.service.Name
	addr := e.service.Addr
	if e.em == nil {
		e.em, _ = endpoints.NewManager(e.cli, name)
	}

	return e.em.AddEndpoint(e.cli.Ctx(), name+"/"+addr, endpoints.Endpoint{Addr: addr})

}

func (e *etcdRegistry) UnRegister() error {
	return e.em.DeleteEndpoint(e.cli.Ctx(), e.service.Name+"/"+e.service.Addr)
}
