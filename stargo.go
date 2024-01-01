package stargo

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/interceptor/recovery"
	"github.com/starfork/stargo/interceptor/validator"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/service"
	"github.com/starfork/stargo/store"
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

	//config *config.Config

	store  map[string]store.Store
	broker broker.Broker

	conf *config.Config
	//mysql  *mysql.Mysql
	//redis  *redis.Redis
	//mongo  *mongo.Mongo
	client   *client.Client
	registry naming.Registry
}

func New(opt ...Option) *App {

	opts := DefaultOptions()
	for _, o := range opt {
		o(&opts)
	}

	conf := opts.Config

	if conf.Timezome != "" {
		time.LoadLocation(conf.Timezome)
	} else {
		time.LoadLocation("Asia/Shanghai")
	}

	s := newServer(opts)

	//注册reflection
	if conf.Environment != ENV_PRODUCTION {
		//app.logger.Debug("env:" + conf.Environment)
		reflection.Register(s)
	}

	app := &App{
		opts:   opts,
		server: s,
		logger: logger.NewZapSugar(conf.Log),
		conf:   conf,
		store:  make(map[string]store.Store),
		//config: opts.Config,
	}

	//注册registry
	if conf.Registry != nil {
		app.conf.Registry.Org = opts.Org
		r := naming.NewRegistry(conf.Registry)
		if err := r.Register(app.Service()); err != nil {
			panic(err)
		}
		app.registry = r
	}
	if conf.Broker != nil {
		//app.broker=
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
		s.logger.Fatalf("failed to listen: %v", err)
	}
	s.logger.Debugf("Starting: gRPC Listener [:%s]\n", port)

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
		log.Fatalf("failed to serve: %v", err)
	}

}

// Stop server
func (s *App) Stop() {

	if s.registry != nil {
		s.logger.Debugf("UnRegister: [%s]\n", s.opts.Name)
		s.registry.UnRegister(s.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	s.server.Stop()
}

// 返回标准服务格式
func (s *App) Service() service.Service {
	return service.Service{
		Org:  s.opts.Org,
		Name: s.opts.Name,
		Addr: s.conf.Port,
	}
}

func (s *App) RegisterService(sd *grpc.ServiceDesc, ss any) *App {
	s.server.RegisterService(sd, ss)
	return s
}

// Restart server
func (s *App) Restart() {
	//mysql 那些重连？
	s.server.GracefulStop()
	s.server.Serve(s.lis)
}

// Server
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
