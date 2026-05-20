package naming

type Registry interface {
	Scheme() string

	Register(service Service) error
	Deregister(service Service) error
	List(name string) []Service
}

var registryFactories = make(map[string]func(*Config) (Registry, error))

func RegisterRegistry(name string, factory func(*Config) (Registry, error)) {
	registryFactories[name] = factory
}

func NewRegistry(name string, conf *Config) (Registry, error) {
	if f, ok := registryFactories[name]; ok {
		return f(conf)
	}
	return nil, nil
}
