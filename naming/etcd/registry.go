package etcd

import (
	"fmt"

	"github.com/starfork/stargo/logger"
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

func NewRegistry(conf *naming.Config) (naming.Registry, error) {
	if conf.Scheme != Scheme {

	}
	cli, err := newClient(conf)
	if err != nil {
		return nil, err
	}

	em, err := endpoints.NewManager(cli, conf.Org)
	if err != nil {
		return nil, err
	}

	return &Registry{
		cli:  cli,
		name: Scheme,
		em:   em,
		conf: conf,
	}, nil
}
func (e *Registry) key(svc naming.Service) string {
	key := e.conf.Org + "/" + svc.Name + "/" + svc.Addr
	fmt.Println(key)
	return key
}

func (e *Registry) Register(svc naming.Service) error {

	p := endpoints.Endpoint{
		Addr: svc.Addr,
	}
	opts := []clientv3.OpOption{}
	if e.conf.Ttl > 0 {
		lease, err := e.cli.Grant(e.cli.Ctx(), e.conf.Ttl)
		if err != nil {
			logger.DefaultLogger.Debugf("waring ttl zero")
		} else {
			opts = append(opts, clientv3.WithLease(lease.ID))
		}
	}
	return e.em.AddEndpoint(e.cli.Ctx(), e.key(svc), p, opts...)
}

func (e *Registry) Deregister(svc naming.Service) error {

	return e.em.DeleteEndpoint(e.cli.Ctx(), e.key(svc))
}

func (e *Registry) List(name string) []naming.Service {

	return nil
}

func (e *Registry) Scheme() string {
	return e.name
}
