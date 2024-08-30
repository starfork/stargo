package server

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App App
type Server struct {
	opts      *Options
	rpcServer *grpc.Server
	lis       net.Listener
	logger    logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry

	conf *Config
	//client *client.Client
}

func New(opt ...Option) *Server {

	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	conf := opts.Config
	time.LoadLocation(opts.Timezone)
	conf.Timezome = opts.Timezone

	s := newRpcServer(opts)

	app := &Server{
		opts:      opts,
		rpcServer: s.(*grpc.Server),
		logger:    logger.DefaultLogger,
		conf:      conf,
		store:     make(map[string]store.Store),
	}

	//注册reflection
	if conf.Env != ENV_PRODUCTION {
		app.logger.Debugf("env:" + conf.Env)
		reflection.Register(app.rpcServer)
	}

	// for k, v := range conf.Store {
	// 	app.Store(k, v)
	// }

	return app
}

// Run   server
func (s *Server) Run() {

	ports := strings.Split(s.conf.RpcServer.Host, ":")
	port := ports[0]
	if len(ports) > 1 {
		port = ports[1] //centos docker 监听ip:port模式有问题
	}
	lis, err := net.Listen("tcp", ":"+port)
	s.lis = lis

	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}
	s.logger.Infof("starting: gRPC Listener %s\n", port)

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

	if err := s.rpcServer.Serve(lis); err != nil {
		s.logger.Fatalf("failed to serve: %v", err)
	}

}

// Stop server
func (s *Server) Stop() {
	s.stopStargo()
	s.rpcServer.Stop()
}
func (s *Server) stopStargo() {
	if s.registry != nil {
		s.logger.Fatalf("UnRegister: [%s]\n", s.opts.Name)
		//s.registry.UnRegister(s.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	if s.broker != nil {
		s.broker.UnSubscribe()
	}
}

// Restart server
func (s *Server) Restart() {
	s.stopStargo()
	s.rpcServer.GracefulStop()
	s.rpcServer.Serve(s.lis)
}

// newServer return new server
func newRpcServer(options *Options) (s grpc.ServiceRegistrar) {

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
