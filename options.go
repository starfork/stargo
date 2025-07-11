package stargo

import (
	"github.com/starfork/stargo/server"

	"google.golang.org/grpc"
)

// Options 参数
type Options struct {
	Addr     string
	Server   *server.Config
	Timezone string

	ServiceDesc []*grpc.ServiceDesc
}

// Option Option
type Option func(o *Options)

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

// UnaryInterceptor Unary server interceptor
func WithUnaryInterceptor(opt grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		o.Server.UnaryInterceptor = append(o.Server.UnaryInterceptor, opt)
	}
}

// StreamInterceptor Stream server interceptor
func WithStreamnIterceptor(opt grpc.StreamServerInterceptor) Option {
	return func(o *Options) {
		o.Server.StreamInterceptor = append(o.Server.StreamInterceptor, opt)
	}
}

// Server option
func WithServer(opt ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.Server.ServerOpts = append(o.Server.ServerOpts, opt...)
	}
}
func WithServerDesc(opt ...*grpc.ServiceDesc) Option {
	return func(o *Options) {
		o.ServiceDesc = append(o.ServiceDesc, opt...)
	}
}

// DefaultOptions default options
func DefaultOptions() *Options {
	o := &Options{
		Server:   server.DefaultConfig,
		Timezone: "Asia/Shanghai",
	}
	return o
}
