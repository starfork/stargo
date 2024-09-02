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
	app.server = server.New()
	app.Init()
	return app
}
func (s *App) Init() {
	conf := s.opts.Config
	time.LoadLocation(s.opts.Timezone)
	conf.Timezome = s.opts.Timezone
	s.once.Do(func() {
		s.logger = logger.DefaultLogger

		r := conf.Registry
		if r != nil {
			r.Environment = conf.Env
			r.Org = conf.Server.Org
			if r.Name == "etcd" {
				s.registry = etcd.NewRegistry(r)
				s.resolver = etcd.NewResolver(r)
			} else if r.Name == "redis" {
				s.registry = redis.NewRegistry(r)
				s.resolver = redis.NewResolver(r)
			} else {
				s.logger.Fatalf("unknow registry")
			}
			s.registry.Register(s.server.Service())
		}
		//注册reflection
		if conf.Env != config.ENV_PRODUCTION {
			s.logger.Debugf("env:" + conf.Env)
			reflection.Register(s.server.Server())
		}

		// for k, v := range conf.Store {
		// 	app.Store(k, v)
		// }

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
