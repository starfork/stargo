package server

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
)

// App App
type Server struct {
	//opts      *Options
	rpcServer *grpc.Server
	lis       net.Listener
	logger    logger.Logger

	conf *Config
}

func New(conf *Config) *Server {

	// opts := DefaultOptions()
	// for _, o := range opt {
	// 	o(opts)
	// }

	//conf := opte.Config

	//time.LoadLocation(opte.Timezone)
	//conf.Timezome = opte.Timezone

	s := newRpcServer(conf)

	app := &Server{

		rpcServer: s.(*grpc.Server),
		logger:    logger.DefaultLogger,
		conf:      conf,
		//store: make(map[string]store.Store),
	}

	return app
}

// Run   server
func (e *Server) Run() {

	ports := strings.Split(e.conf.Addr, ":")
	port := ports[0]
	if len(ports) > 1 {
		port = ports[1] //centos docker 监听ip:port模式有问题
	}
	lis, err := net.Listen("tcp", ":"+port)
	e.lis = lis

	if err != nil {
		e.logger.Fatalf("failed to listen: %v", err)
	}

	e.logger.Infof("starting: gRPC Listener %s\n", e.conf.Addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		e.Stop()

		if i, ok := sg.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	if err := e.rpcServer.Serve(lis); err != nil {
		e.logger.Fatalf("failed to serve: %v", err)
	}

}

// Stop server
func (e *Server) Stop() {
	e.rpcServer.Stop()
}

// Restart server
func (e *Server) Restart() {

	e.rpcServer.GracefulStop()
	e.rpcServer.Serve(e.lis)
}

// newServer return new server
func newRpcServer(conf *Config) (s grpc.ServiceRegistrar) {

	opt := append(conf.Server,
		grpc.ChainUnaryInterceptor(conf.UnaryInterceptor...),
		grpc.ChainStreamInterceptor(conf.StreamInterceptor...),
	)

	// if conf.Xds {
	// 	var err error
	// 	if s, err = xde.NewGRPCServer(opt...); err != nil {
	// 		panic(err)
	// 	}
	// } else {
	s = grpc.NewServer(opt...)
	//}

	return s
}

func (e *Server) Service() naming.Service {
	return naming.Service{
		Org:  e.conf.Org,
		Name: e.conf.Name,
		Addr: e.conf.Addr,
	}
}

func (e *Server) Server() *grpc.Server {
	return e.rpcServer
}
