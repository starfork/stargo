package stargo

import (
	"github.com/starfork/stargo/cache"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/debug/tracer"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"

	"google.golang.org/grpc"
)

// Options 参数
type Options struct {
	Name             string
	Config           *config.Config
	UnaryInterceptor []grpc.UnaryServerInterceptor

	Server   []grpc.ServerOption
	Registry naming.Registry
	Tracer   tracer.Tracer
	Logger   logger.Logger
	Cache    cache.Cache
}

// Option Option
type Option func(o *Options)

func Name(c string) Option {
	return func(o *Options) {
		o.Name = c
	}
}
func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// UnaryInterceptor Unary server interceptor
func UnaryInterceptor(opt ...grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.UnaryInterceptor = append(o.UnaryInterceptor, opt...)
	}
}

// Server option
func Server(opt ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.Server = append(o.Server, opt...)
	}
}

// DefaultOptions default options
func DefaultOptions() Options {
	o := Options{
		Name: "strago",
		// Conf: &config.Config{
		// 	Deploy: "Monolithic",
		// 	Base:   &config.ServerConfig{},
		// },
		//Name:     "Default",

		//LogFile:  "debug.log",
	}
	return o
}
