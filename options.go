package strago

import "google.golang.org/grpc"

//Options 参数
type Options struct {
	Name string
	Port string
	//Server           *grpc.Server
	UnaryInterceptor []grpc.UnaryServerInterceptor
	Registry         string
	Reflect          bool
	Pb               interface{}
}

//Option Option
type Option func(o *Options)

//Name set server name
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

//Port set server name
func Port(port string) Option {
	return func(o *Options) {
		o.Port = port
	}
}

//Reflect set server name
func Reflect() Option {
	return func(o *Options) {
		o.Reflect = true
	}
}

//Pb set server name
func Pb(pb func()) Option {
	return func(o *Options) {
		o.Pb = pb
	}
}

//UnaryInterceptor Unary server interceptor
func UnaryInterceptor(opt ...grpc.UnaryServerInterceptor) Option {
	return func(o *Options) {
		for _, v := range opt {
			o.UnaryInterceptor = append(o.UnaryInterceptor, v)
		}
	}
}

//DefaultOptions default options
func DefaultOptions() Options {
	o := Options{
		Name:     "Default",
		Balancer: "",
	}
	return o
}
