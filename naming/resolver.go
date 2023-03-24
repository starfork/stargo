package naming

type Resolver interface {
	Register(service Service) error
	UnRegister() error
}
