package stargo

import "github.com/starfork/stargo/server"

type Options struct {
	Name, Port string
	serverOpts server.Options
	//balancerOpts
	//brokerOpts
	Server server.Server
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.serverOpts.Name = n
	}
}

func Port(p string) Option {
	return func(o *Options) {
		o.serverOpts.Port = p
	}
}
