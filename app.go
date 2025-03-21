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
	ctx context.Context

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

func New(opt ...Option) *App {

	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}
	s := &App{
		ctx:   context.Background(),
		opts:  opts,
		store: make(map[string]store.Store),
	}
	s.Init()

	return s
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

		for k, v := range conf.Store {
			if k == "mysql" {
				s.Store(k, smysql.NewMysql(v))
			}
			if k == "redis" {
				s.Store(k, sredis.NewRedis(v))
			}
		}
		if conf.Broker != nil {
			s.broker = nats.NewBroker(conf.Broker)
		}
		if conf.Tracer != nil {
			var err error
			if s.tracer, err = otel.NewTracer(conf.Tracer); err != nil {
				s.logger.Fatalf("tracer init fail: [%s]\n", conf.Tracer.Host)
			}
			//
			s.opts.Config.Server.Server = append(
				s.opts.Config.Server.Server,
				grpc.StatsHandler(otelgrpc.NewServerHandler()),
			)
		}

		if conf.Registry != nil {
			r := conf.Registry
			var err error
			if r.Scheme == "etcd" {
				if s.registry, err = etcd.NewRegistry(r); err != nil {
					s.logger.Fatalf("unknow registry")
				}

				if s.resolver, err = etcd.NewResolver(r); err != nil {
					s.logger.Fatalf("unknow registry")
				}
			} else {
				s.logger.Fatalf("unknow registry")
			}
		}
	})
}

// Run   server
func (s *App) Run() {
	conf := s.opts.Config
	s.server = server.New(conf.Server)

	//注册reflection
	if conf.Env != config.ENV_PRODUCTION {
		reflection.Register(s.server.Server())
	}
	service := s.server.Service()
	service.Org = conf.Registry.Org
	if err := s.registry.Register(service); err != nil {
		s.logger.Fatalf("registry err %+v", err)
	}
	s.server.Run()
}
func (s *App) RunService(desc *grpc.ServiceDesc, impl any) {
	conf := s.opts.Config
	s.server = server.New(conf.Server)

	//注册reflection
	if conf.Env != config.ENV_PRODUCTION {
		reflection.Register(s.server.Server())
	}
	service := s.server.Service()
	service.Org = conf.Registry.Org
	if err := s.registry.Register(service); err != nil {
		s.logger.Fatalf("registry err %+v", err)
	}

	s.server.Server().RegisterService(desc, impl)
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

	if s.tracer != nil {
		s.tracer.Close()
	}
}

// Restart server
func (s *App) Restart() {
	s.stopStargo()
	s.server.Restart()
	s.Init()
}
