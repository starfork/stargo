package naming

type Registry interface {
	Register(service Service) error
	UnRegister() error
}
