package stargo

import (
	"context"
	"sync"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/broker/nats"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/naming/etcd"
	"github.com/starfork/stargo/server"
	"github.com/starfork/stargo/store"
	smysql "github.com/starfork/stargo/store/mysql"
	sredis "github.com/starfork/stargo/store/redis"
	"github.com/starfork/stargo/tracer"
	"github.com/starfork/stargo/tracer/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// App App
type App struct {
	ctx  context.Context
	name string

	conf   *config.Config
	opts   *Options
	server *server.Server
	logger logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry
	resolver naming.Resolver
	tracer   tracer.Tracer

	Tz *time.Location

	once sync.Once
}

func New(name string, conf *config.Config) *App {

	opts := DefaultOptions()

	s := &App{
		ctx:   context.Background(),
		opts:  opts,
		store: make(map[string]store.Store),
		name:  name,
		conf:  conf,
	}
	s.initConfig()
	return s
}

// init by Config
func (s *App) initConfig() {

	s.once.Do(func() {
		s.logger = logger.DefaultLogger

		for k, v := range s.conf.Store {
			if k == "mysql" {
				s.Store(k, smysql.NewMysql(v))
			}
			if k == "redis" {
				s.Store(k, sredis.NewRedis(v))
			}
		}
		if s.conf.Broker != nil {
			s.conf.Broker.App = s.name
			s.broker = nats.NewBroker(s.conf.Broker)
		}
		if s.conf.Tracer != nil {
			var err error
			if s.tracer, err = otel.NewTracer(s.conf.Tracer); err != nil {
				s.logger.Fatalf("tracer init fail: [%s]\n", s.conf.Tracer.Host)
			}
			//
			s.conf.Server.ServerOpts = append(
				s.conf.Server.ServerOpts,
				grpc.StatsHandler(otelgrpc.NewServerHandler()),
			)
		}

		if s.conf.Registry != nil {
			r := s.conf.Registry

			var err error
			if r.Scheme == "etcd" {
				if s.registry, err = etcd.NewRegistry(r); err != nil {
					s.logger.Fatalf("etcd registry err %+v", err.Error())
				}

				if s.resolver, err = etcd.NewResolver(r); err != nil {
					s.logger.Fatalf("etcd resolver %+v", err)
				}
			} else {
				s.logger.Fatalf("unknow registry")
			}
		}
	})

}

// 初始化数据库之类的东西
func (s *App) Init(opt ...Option) {

	for _, o := range opt {
		o(s.opts)
	}
	//其他的需要传递给
}

// Run   server

func (s *App) beforeRun() {
	if tz, err := time.LoadLocation(s.opts.Timezone); err == nil {
		s.Tz = tz
		s.conf.Timezome = s.opts.Timezone
	}

	sc := s.opts.Server
	sc.Addr = s.conf.Server.Addr
	sc.Name = s.name
	s.server = server.New(sc)

	//注册reflection
	if s.conf.Env != config.ENV_PRODUCTION {
		reflection.Register(s.server.Server())
	}
	service := s.server.Service()

	if err := s.registry.Register(service); err != nil {
		s.logger.Fatalf("registry err %+v", err)
	}

}

func (s *App) Run(desc *grpc.ServiceDesc, impl any) {
	s.beforeRun()
	s.server.Server().RegisterService(desc, impl)

	for _, v := range s.opts.ServiceDesc {
		s.server.Server().RegisterService(v, impl)
	}
	s.server.Run()
}

// Stop server
func (s *App) Stop() {
	s.stopStargo()
	s.server.Stop()
}
func (s *App) stopStargo() {
	if s.registry != nil {
		s.logger.Infof("unregister: [%s]\n", s.name)
		s.registry.Deregister(s.server.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	if s.broker != nil {
		s.broker.UnSubscribe()
	}

	if s.tracer != nil {
		s.tracer.Close()
	}
}

// Restart server
func (s *App) Restart() {
	s.stopStargo()
	s.server.Restart()
	s.initConfig()
	s.Init()
}
