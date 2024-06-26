package stargo

import (
	"github.com/starfork/stargo/config"

	"google.golang.org/grpc"
)

// Options 参数
type Options struct {
	Org, Name, Addr   string
	Config            *config.Config
	UnaryInterceptor  []grpc.UnaryServerInterceptor
	StreamInterceptor []grpc.StreamServerInterceptor

	Server []grpc.ServerOption
	//Registry naming.Registry
	//Tracer tracer.Tracer
	//Logger logger.Logger
	//Cache  cache.Cache
}

// Option Option
type Option func(o *Options)

func WithOrg(c string) Option {
	return func(o *Options) {
		o.Org = c
	}
}
func WithName(c string) Option {
	return func(o *Options) {
		o.Name = c
	}
}

func WithAddr(c string) Option {
	return func(o *Options) {
		o.Addr = c
	}
}
func WithConfig(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// UnaryInterceptor Unary server interceptor
func WithUnaryInterceptor(opt grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.UnaryInterceptor = append(o.UnaryInterceptor, opt)
	}
}

// StreamInterceptor Stream server interceptor
func WithStreamnIterceptor(opt grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.StreamInterceptor = append(o.StreamInterceptor, opt)
	}
}

// Server option
func WithServer(opt ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.Server = append(o.Server, opt...)
	}
}

// DefaultOptions default options
func DefaultOptions() *Options {
	o := &Options{
		Org:  "stargo", //与proto文件对称即可。比如stargo.service.ServiceHandler
		Name: "service",
	}
	return o
}
