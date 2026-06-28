package stargo

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/server"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/tracer"
	"github.com/starfork/stargo/util/ustring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// App App
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	name   string

	conf   *config.Config
	opts   *Options
	server *server.Server
	logger logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry
	resolver naming.Resolver
	tracer   tracer.Tracer

	once sync.Once
}

func New(name string, conf *config.Config) *App {

	opts := DefaultOptions()
	ctx, cancel := context.WithCancel(context.Background())

	s := &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   opts,
		store:  make(map[string]store.Store),
		name:   name,
		conf:   conf,
	}
	s.initConfig()
	return s
}

// init by Config
func (s *App) initConfig() {

	s.once.Do(func() {
		s.logger = logger.DefaultLogger

		for k, v := range s.conf.Store {
			if st := store.NewStore(k, v); st != nil {
				if k == "mysql" {
					v.TimeLocation = s.opts.Timezone
					if v.Prefix == "" {
						v.Prefix = ustring.Or(s.name+"_", os.Getenv("MYSQL_PREFIX"))
					}
				}
				s.Store(k, st)
			}
		}
		if s.conf.Broker != nil {
			s.conf.Broker.App = s.name
			if b, err := broker.NewBroker(s.conf.Broker.Name, s.conf.Broker); err != nil {
				s.logger.Warnf("broker init error: %v", err)
			} else if b != nil {
				s.broker = b
			}
		}
		s.tracer = tracer.DefaultTracer

		if s.conf.Registry != nil {
			r := s.conf.Registry
			var err error
			if s.registry, err = naming.NewRegistry(r.Scheme, r); err != nil {
				s.logger.Fatalf("registry err %+v", err.Error())
			}
			if s.resolver, err = naming.NewResolver(r.Scheme, r); err != nil {
				s.logger.Fatalf("resolver err %+v", err)
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
		time.Local = tz
		//s.Tz = tz
		s.conf.Timezone = s.opts.Timezone
	}

	sc := s.opts.Server
	sc.Addr = s.conf.Server.Addr
	sc.Name = s.name
	s.server = server.New(sc)
	s.server.Health().SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	//注册reflection
	if s.conf.Env != config.ENV_PRODUCTION {
		reflection.Register(s.server.Server())
	}
	service := s.server.Service()

	if s.registry != nil {
		if err := s.registry.Register(service); err != nil {
			s.logger.Fatalf("registry err %+v", err)
		}
	}

}

func (s *App) Run(desc *grpc.ServiceDesc, impl any) {
	s.beforeRun()
	s.server.Server().RegisterService(desc, impl)

	for _, v := range s.opts.ServiceDesc {
		s.server.Server().RegisterService(v, impl)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		s.logger.Infof("received signal: %v, shutting down", sig)
		s.Stop()
		if i, ok := sig.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	s.server.Run()
}

// Stop server
func (s *App) Stop() {
	s.cancel()
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
