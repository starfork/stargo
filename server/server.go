package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/starfork/strago/interceptor/recovery"
	"github.com/starfork/strago/interceptor/validator"
	"github.com/starfork/strago/interceptor/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//grpc-cli -service Ucenter -method List localhost:50501
)

//App App
type Server struct {
	Option     Options
	GRPCServer *grpc.Server
}

//New app
func New(opts ...Option) *Server {

	time.LoadLocation("Asia/Shanghai")

	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	s := newServer(options.UnaryInterceptor...)

	if options.Reflect {
		reflection.Register(s)
	}

	if options.Registry != "" {
		//log.Printf("Balancer: [%s]\n", options.Balancer)
		//go registry.Register(options.Registry, options.Name, options.Port, 5)
	}

	return &Server{
		Option:     options,
		GRPCServer: s,
	}
	//return s

}

//Run   server
func (s *Server) Run() {
	opt := s.Option
	lis, err := net.Listen("tcp", opt.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting: gRPC Listener [%s]\n", opt.Port)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		fmt.Println(s)
		// if opt.Balancer != "" {
		// 	log.Printf("UnRegister: [%s]\n", opt.Name)
		// 	balancer.UnRegister(opt.Name, opt.Balancer)
		// }

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err := s.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// Stop server
func (s *Server) Stop() {
	//os.Exit(0)
}

// Restart server
func (s *Server) Restart() {

}

//Server set server name
func (s *Server) Server() *grpc.Server {
	return s.GRPCServer
}

//newServer return new server
func newServer(interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptor = append(interceptor,
		validator.Unary(),
		zap.Unary(),
		recovery.Unary(),
	)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				interceptor...,
			),
		),
	)
	return s
}
