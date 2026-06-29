package server

import (
	"net"
	"strings"
	"time"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// App App
type Server struct {
	rpcServer *grpc.Server
	lis       net.Listener
	logger    logger.Logger
	conf      *Config
	health    *HealthServer
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
		health:    NewHealthServer(),
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

	e.rpcServer.Serve(lis)
}

// Stop server
func (e *Server) Stop() {
	done := make(chan struct{})
	go func() {
		e.rpcServer.GracefulStop()
		close(done)
	}()
	timeout := e.conf.ShutdownTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
		e.logger.Warnf("graceful stop timed out after %v, forcing stop", timeout)
		e.rpcServer.Stop()
	}
}

// Restart server
func (e *Server) Restart() {
	e.Stop()
	e.rpcServer.Serve(e.lis)
}

func (e *Server) Health() *HealthServer {
	return e.health
}

func (e *Server) registerHealth() {
	grpc_health_v1.RegisterHealthServer(e.rpcServer, e.health)
}

// newServer return new server
func newRpcServer(conf *Config) (s grpc.ServiceRegistrar) {
	interceptors := conf.UnaryInterceptor
	streamInterceptors := conf.StreamInterceptor

	if conf.Metrics {
		interceptors = append([]grpc.UnaryServerInterceptor{UnaryServerMetricsInterceptor}, interceptors...)
		streamInterceptors = append([]grpc.StreamServerInterceptor{StreamServerMetricsInterceptor}, streamInterceptors...)
	}

	opt := append(conf.ServerOpts,
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	s = grpc.NewServer(opt...)
	return s
}

func (e *Server) Service() naming.Service {
	return naming.Service{
		Name: e.conf.Name,
		Addr: e.conf.Addr,
	}
}

func (e *Server) Server() *grpc.Server {
	return e.rpcServer
}
