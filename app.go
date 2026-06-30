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
	
	// Saved service descriptors for restart
	serviceDesc  *grpc.ServiceDesc
	serviceImpl  any
	extraDescs   []*grpc.ServiceDesc
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
	s.logger = logger.DefaultLogger
	s.initLogger()
	return s
}

// init by Config
func (s *App) initConfig() {

	s.once.Do(func() {
		s.initStore()
		s.initBroker()
		s.initTracer()
		s.initRegistry()
	})

}

func (s *App) initLogger() {
	if s.conf.Log != nil && s.conf.Log.Driver != "" {
		if l, err := logger.NewLogger(s.conf.Log.Driver, s.conf.Log); err == nil {
			s.logger = l
		}
	}
}

func (s *App) initStore() {
	for k, v := range s.conf.Store {
		if st := store.NewStore(k, v); st != nil {
			v.TimeLocation = s.opts.Timezone
			if v.Prefix == "" {
				v.Prefix = s.name + "_"
			}
			s.Store(k, st)
			if s.server != nil {
				s.server.Health().SetDependency("store/"+k, true)
			}
		}
	}
}

func (s *App) initBroker() {
	if s.conf.Broker != nil {
		s.conf.Broker.App = s.name
		if b, err := broker.NewBroker(s.conf.Broker.Name, s.conf.Broker); err != nil {
			s.logger.Warnf("broker init error: %v", err)
		} else if b != nil {
			s.broker = b
			if s.server != nil {
				s.server.Health().SetDependency("broker/"+s.conf.Broker.Name, true)
			}
		}
	}
}

func (s *App) initTracer() {
	if s.conf.Tracer != nil {
		if t, err := tracer.NewTracer(s.conf.Tracer.Driver, s.conf.Tracer); err == nil && t != nil {
			s.tracer = t
		} else {
			s.tracer = tracer.DefaultTracer
		}
	} else {
		s.tracer = tracer.DefaultTracer
	}
}

func (s *App) initRegistry() {
	if s.conf.Registry != nil {
		r := s.conf.Registry
		var err error
		if s.registry, err = naming.NewRegistry(r.Scheme, r); err != nil {
			s.logger.Fatalf("registry err %+v", err.Error())
		}
		if s.resolver, err = naming.NewResolver(r.Scheme, r); err != nil {
			s.logger.Fatalf("resolver err %+v", err)
		}
		if s.server != nil {
			s.server.Health().SetDependency("registry/"+r.Scheme, true)
		}
	}
}

// 初始化数据库之类的东西
func (s *App) Init(opt ...Option) {

	for _, o := range opt {
		o(s.opts)
	}
	s.initConfig()
}

// Run   server

func (s *App) beforeRun() {
	// Ensure config is initialized (idempotent, safe if Init() was already called)
	s.initConfig()
	
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
	// Save service descriptors for restart
	s.serviceDesc = desc
	s.serviceImpl = impl
	s.extraDescs = s.opts.ServiceDesc
	
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
	// 1. Set health to NOT_SERVING to signal load balancers
	if s.server != nil && s.server.Health() != nil {
		s.server.Health().SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	// 2. Drain in-flight requests
	drainTimeout := s.server.Config().ShutdownTimeout / 2
	if drainTimeout < 5*time.Second {
		drainTimeout = 5 * time.Second
	}
	s.logger.Infof("draining in-flight requests for %v", drainTimeout)
	time.Sleep(drainTimeout)

	// 3. Deregister from service registry
	if s.registry != nil {
		if s.server != nil {
			s.logger.Infof("unregister: [%s]", s.name)
			s.registry.Deregister(s.server.Service())
		}
		s.registry.Close()
	}

	if s.resolver != nil {
		s.resolver.Close()
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
	// Reinitialize dependencies
	s.initLogger()
	s.initStore()
	s.initBroker()
	s.initTracer()
	s.initRegistry()
	// Recreate and start server
	s.server.Restart()
	
	// Re-register business services
	if s.serviceDesc != nil && s.serviceImpl != nil {
		s.server.Server().RegisterService(s.serviceDesc, s.serviceImpl)
		for _, v := range s.extraDescs {
			s.server.Server().RegisterService(v, s.serviceImpl)
		}
	}
	
	// Re-register to registry
	service := s.server.Service()
	if s.registry != nil {
		if err := s.registry.Register(service); err != nil {
			s.logger.Errorf("registry err %+v", err)
		}
	}
	
	s.server.Run()
}
