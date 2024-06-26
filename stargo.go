package stargo

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ENV_DEV        = "dev"        //本地测试环境
	ENV_DOCKER     = "docker"     //docker模式
	ENV_PRODUCTION = "production" //正式环境
)

// App App
type App struct {
	opts   *Options
	server *grpc.Server
	lis    net.Listener
	logger logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry

	conf *config.Config
	//client *client.Client
}

func New(opt ...Option) *App {

	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	conf := opts.Config
	tz := "Asia/Shanghai"
	if conf.Timezome != "" {
		tz = conf.Timezome
	}
	time.LoadLocation(tz)

	s := newServer(opts)

	app := &App{
		opts:   opts,
		server: s.(*grpc.Server),
		logger: logger.DefaultLogger,
		conf:   conf,
		store:  make(map[string]store.Store),
	}

	//注册reflection
	if conf.Environment != ENV_PRODUCTION {
		app.logger.Debugf("env:" + conf.Environment)
		reflection.Register(app.server)
	}

	return app
}

// Run   server
func (s *App) Run() {

	//	s.logger.Debugf("ServerPort%+v", s.conf.ServerPort)
	ports := strings.Split(s.conf.Port, ":")
	port := ports[0]
	if len(ports) > 1 {
		port = ports[1] //centos docker 监听ip:port模式有问题
	}
	lis, err := net.Listen("tcp", ":"+port)
	s.lis = lis

	if err != nil {
		s.logger.Logf(logger.FatalLevel, "failed to listen: %v", err)
	}
	s.logger.Logf(logger.DebugLevel, "starting: gRPC Listener %s\n", port)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		s.Stop()

		if i, ok := sg.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	if err := s.server.Serve(lis); err != nil {
		s.logger.Logf(logger.FatalLevel, "failed to serve: %v", err)
	}

}

// Stop server
func (s *App) Stop() {
	s.stopStargo()
	s.server.Stop()
}
func (s *App) stopStargo() {
	if s.registry != nil {
		s.logger.Logf(logger.FatalLevel, "UnRegister: [%s]\n", s.opts.Name)
		s.registry.UnRegister(s.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	if s.broker != nil {
		s.broker.UnSubscribe()
	}
}

func (s *App) RegisterService(sd *grpc.ServiceDesc, ss any) *App {
	s.server.RegisterService(sd, ss)
	return s
}

// Restart server
func (s *App) Restart() {
	s.stopStargo()
	s.server.GracefulStop()
	s.server.Serve(s.lis)
}

// newServer return new server
func newServer(options *Options) (s grpc.ServiceRegistrar) {

	opt := append(options.Server,
		grpc.ChainUnaryInterceptor(options.UnaryInterceptor...),
		grpc.ChainStreamInterceptor(options.StreamInterceptor...),
	)

	// if conf.Xds {
	// 	var err error
	// 	if s, err = xds.NewGRPCServer(opt...); err != nil {
	// 		panic(err)
	// 	}
	// } else {
	s = grpc.NewServer(opt...)
	//}

	return s
}
