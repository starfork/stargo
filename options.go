package stargo

import (
	"github.com/starfork/stargo/config"

	"google.golang.org/grpc"
)

// Options 参数
type Options struct {
	Org, Name, Addr  string
	Config           *config.Config
	UnaryInterceptor []grpc.UnaryServerInterceptor

	Server []grpc.ServerOption
	//Registry naming.Registry
	//Tracer tracer.Tracer
	//Logger logger.Logger
	//Cache  cache.Cache
}

// Option Option
type Option func(o *Options)

func Org(c string) Option {
	return func(o *Options) {
		o.Org = c
	}
}
func Name(c string) Option {
	return func(o *Options) {
		o.Name = c
	}
}

func Addr(c string) Option {
	return func(o *Options) {
		o.Addr = c
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
		Org:  "stargo", //与proto文件对称即可。比如stargo.service.ServiceHandler
		Name: "service",
	}
	return o
}
