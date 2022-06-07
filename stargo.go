package stargo

import (
	"github.com/starfork/stargo/server"
	"google.golang.org/grpc"
)

type App interface {
	Server() server.Server
}

type Option func(*Options)

type app struct {
	opts Options
}

func New(opts ...Option) *app {
	//server.New(opts.)
	return &app{}
}

func (s *app) Server() *grpc.Server {
	return s.opts.Server.GRPCServer
}
func (s *app) Run() error {
	return nil
}
