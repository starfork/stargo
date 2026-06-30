package server

import (
	"net"
	"time"

	"github.com/starfork/stargo/interceptor/recovery"
	"github.com/starfork/stargo/interceptor/timeout"
	internaltls "github.com/starfork/stargo/internal/tls"
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
	app.registerHealth()

	return app
}

// Run   server
func (e *Server) Run() {

	_, port, err := net.SplitHostPort(e.conf.Addr)
	if err != nil {
		// If parsing fails, assume it's just a port
		port = e.conf.Addr
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
	// Create new listener
	_, port, err := net.SplitHostPort(e.conf.Addr)
	if err != nil {
		// If parsing fails, assume it's just a port
		port = e.conf.Addr
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		e.logger.Fatalf("failed to listen: %v", err)
	}
	e.lis = lis
	
	// Create new RPC server
	s := newRpcServer(e.conf)
	e.rpcServer = s.(*grpc.Server)
	e.health = NewHealthServer()
	e.registerHealth()
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

	// Add default timeout interceptor if configured
	if conf.DefaultTimeout > 0 {
		interceptors = append([]grpc.UnaryServerInterceptor{timeout.Unary(conf.DefaultTimeout)}, interceptors...)
		streamInterceptors = append([]grpc.StreamServerInterceptor{timeout.Stream(conf.DefaultTimeout)}, streamInterceptors...)
	}

	// Add recovery interceptors by default
	interceptors = append([]grpc.UnaryServerInterceptor{recovery.Unary()}, interceptors...)
	streamInterceptors = append([]grpc.StreamServerInterceptor{recovery.Stream()}, streamInterceptors...)

	if conf.Metrics {
		interceptors = append([]grpc.UnaryServerInterceptor{UnaryServerMetricsInterceptor}, interceptors...)
		streamInterceptors = append([]grpc.StreamServerInterceptor{StreamServerMetricsInterceptor}, streamInterceptors...)
	}

	opt := append(conf.ServerOpts,
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// Add TLS if configured
	if conf.CertFile != "" && conf.KeyFile != "" {
		creds, err := internaltls.NewServerTransportCredentials(conf.CertFile, conf.KeyFile, conf.CAFile)
		if err == nil {
			opt = append(opt, grpc.Creds(creds))
		}
	}

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

func (e *Server) Config() *Config {
	return e.conf
}
