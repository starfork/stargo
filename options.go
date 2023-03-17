package stargo

import (
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/registry"

	"google.golang.org/grpc"
)

// Options 参数
type Options struct {
	Balancer         string
	Conf             *config.Config
	UnaryInterceptor []grpc.UnaryServerInterceptor

	Server   []grpc.ServerOption
	Registry registry.Registry
}

// Option Option
type Option func(o *Options)

func Conf(c *config.Config) Option {
	return func(o *Options) {
		o.Conf = c
	}
}

// Balancer set server name
func Balancer(b string) Option {
	return func(o *Options) {
		o.Balancer = b
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
		// Conf: &config.Config{
		// 	Deploy: "Monolithic",
		// 	Base:   &config.ServerConfig{},
		// },
		//Name:     "Default",
		Balancer: "",
		//LogFile:  "debug.log",
	}
	return o
}
