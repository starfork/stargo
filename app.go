package stargo

import (
	"context"
	"sync"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/naming/etcd"
	"github.com/starfork/stargo/naming/redis"
	"github.com/starfork/stargo/server"
	"github.com/starfork/stargo/store"
	smysql "github.com/starfork/stargo/store/mysql"
	sredis "github.com/starfork/stargo/store/redis"
	"google.golang.org/grpc/reflection"
)

// App App
type App struct {
	ctx context.Context

	opts   *Options
	server *server.Server
	logger logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry
	resolver naming.Resolver

	Tz *time.Location

	once sync.Once
	//client *client.Client
}

func New(opt ...Option) *App {

	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}
	app := &App{
		ctx:   context.Background(),
		opts:  opts,
		store: make(map[string]store.Store),
	}

	app.server = server.New(opts.Config.Server)
	app.Init()
	return app
}
func (s *App) Init() {
	conf := s.opts.Config
	if tz, err := time.LoadLocation(s.opts.Timezone); err == nil {
		s.Tz = tz
		conf.Timezome = s.opts.Timezone
	}
	//registry，store这样的方式，需要改进成配置形式
	s.once.Do(func() {
		s.logger = logger.DefaultLogger

		r := conf.Registry
		if r != nil {
			if r.Scheme == "etcd" {
				rg, err := etcd.NewRegistry(r)
				if err != nil {
					s.logger.Fatalf("unknow registry")
				}
				s.registry = rg
				rs, err := etcd.NewResolver(r)
				if err != nil {
					s.logger.Fatalf("unknow registry")
				}
				s.resolver = rs
			} else if r.Scheme == "redis" {
				s.registry = redis.NewRegistry(r)
				s.resolver = redis.NewResolver(r)
			} else {
				s.logger.Fatalf("unknow registry")
			}
			service := s.server.Service()
			service.Org = r.Org
			err := s.registry.Register(service)
			if err != nil {
				s.logger.Fatalf("registry err %+v", err)
			}
		}
		//注册reflection
		if conf.Env != config.ENV_PRODUCTION {
			reflection.Register(s.server.Server())
		}

		for k, v := range conf.Store {
			if k == "mysql" {
				s.Store(k, smysql.NewMysql(v))
			}
			if k == "redis" {
				s.Store(k, sredis.NewRedis(v))
			}
		}

	})
}

// Run   server
func (s *App) Run() {
	s.server.Run()
}

// Stop server
func (s *App) Stop() {
	s.stopStargo()
	s.server.Stop()
}
func (s *App) stopStargo() {
	if s.registry != nil {
		s.logger.Infof("UnRegister: [%s]\n", s.opts.Name)
		s.registry.Deregister(s.server.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	if s.broker != nil {
		s.broker.UnSubscribe()
	}
}

// Restart server
func (s *App) Restart() {
	s.stopStargo()
	s.server.Restart()
	s.Init()
}
