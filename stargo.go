package stargo

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/interceptor/recovery"
	"github.com/starfork/stargo/interceptor/validator"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/service"
	"github.com/starfork/stargo/store/mysql"
	"github.com/starfork/stargo/store/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	sf "github.com/sony/sonyflake"
)

var (
	ENV_DEV        = "dev"        //本地测试环境
	ENV_DOCKER     = "docker"     //docker模式
	ENV_PRODUCTION = "production" //正式环境
)

// App App
type App struct {
	opts   Options
	server *grpc.Server
	lis    net.Listener
	logger *zap.SugaredLogger
	sfid   *sf.Sonyflake

	config *config.Config
	conf   *config.ServerConfig
	mysql  *mysql.Mysql
	redis  *redis.Redis

	registry naming.Registry
}

func NewApp(opts ...Option) *App {
	return New(opts...)
}

func New(opts ...Option) *App {

	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	conf := options.Config.GetServerConfig()

	if conf.Timezome != "" {
		time.LoadLocation(conf.Timezome)
	} else {
		time.LoadLocation("Asia/Shanghai")
	}

	if conf.ServerName != "" {
		options.Name = conf.ServerName
	}

	s := newServer(options)

	app := &App{
		opts:   options,
		server: s,
		logger: logger.NewZapSugar(conf.Log),
		conf:   conf,
		config: options.Config,
		//registry: conf.Registry.Name,

		//	Loger:  log.Sugar,
	}
	//注册reflection
	if conf.Environment != ENV_PRODUCTION {
		app.logger.Debug("env:" + conf.Environment)
		reflection.Register(s)
	}
	//注册registry
	if conf.Registry != nil {

		r := naming.NewRegistry(conf.Registry)
		//options.Name, options.Port, 1800
		if err := r.Register(service.Service{
			Name: conf.ServerName,
			Addr: conf.ServerPort,
		}); err != nil {
			panic(err)
		}

		app.registry = r
	}

	return app
}

// Run   server
func (s *App) Run() {

	s.logger.Debugf("ServerPort%+v", s.conf.ServerPort)
	lis, err := net.Listen("tcp", s.conf.ServerPort)
	s.lis = lis

	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}
	s.logger.Debugf("Starting: gRPC Listener [%s]\n", s.conf.ServerPort)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		s.server.Stop()
		if i, ok := sg.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// Stop server
func (s *App) Stop() {

	if s.registry != nil {
		s.logger.Debugf("UnRegister: [%s]\n", s.conf.ServerName)
		//s.registry.UnRegister(s.Service())
	}

	s.server.Stop()
}

// 返回标准服务格式
func (s *App) Service() service.Service {
	return service.Service{
		Name: s.opts.Name,
		Addr: s.opts.Addr,
	}
}

// Restart server
func (s *App) Restart() {
	s.server.GracefulStop()
	s.server.Serve(s.lis)
}

// Server set server name
func (s *App) Server() *grpc.Server {
	return s.server
}

// newServer return new server
func newServer(options Options) *grpc.Server {
	//var opt []grpc.ServerOption
	//目前只测试了unaryserver
	opt := append(options.Server, interceptors(options.UnaryInterceptor...))

	//grpc.StatsHandler()
	s := grpc.NewServer(opt...)

	return s
}

func interceptors(interceptor ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	interceptor = append(interceptor,
		validator.Unary(),
		//zap.Unary(),
		recovery.Unary(),
	)
	opt := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			interceptor...,
		),
	)
	return opt
}
