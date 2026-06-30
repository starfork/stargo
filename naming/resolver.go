package naming

import "google.golang.org/grpc/resolver"

type Resolver interface {
	resolver.Builder
	Config() *Config
	Close() error
}

var resolverFactories = make(map[string]func(*Config) (Resolver, error))

func RegisterResolver(name string, factory func(*Config) (Resolver, error)) {
	resolverFactories[name] = factory
}

func NewResolver(name string, conf *Config) (Resolver, error) {
	if f, ok := resolverFactories[name]; ok {
		return f(conf)
	}
	return nil, nil
}
