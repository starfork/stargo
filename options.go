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
	Server            []grpc.ServerOption
	Timezone          string
}

// Option Option
type Option func(o *Options)

func WithOrg(c string) Option {
	return func(o *Options) {
		o.Config.Registry.Org = c
	}
}
func WithName(c string) Option {
	return func(o *Options) {
		o.Name = c
	}
}
func WithTimezome(c string) Option {
	return func(o *Options) {
		o.Timezone = c
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
		o.Config.Server.UnaryInterceptor = append(o.Config.Server.UnaryInterceptor, opt)
	}
}

// StreamInterceptor Stream server interceptor
func WithStreamnIterceptor(opt grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.Config.Server.StreamInterceptor = append(o.Config.Server.StreamInterceptor, opt)
	}
}

// Server option
func WithServer(opt ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.Config.Server.Server = append(o.Config.Server.Server, opt...)
	}
}

// DefaultOptions default options
func DefaultOptions() *Options {
	o := &Options{
		Config:   config.DefaultConfig,
		Org:      "stargo", //与proto文件对称即可。比如stargo.service.ServiceHandler
		Name:     "service",
		Timezone: "Asia/Shanghai",
	}
	return o
}
