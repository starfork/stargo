package naming

type Registry interface {
	Scheme() string

	Register(service Service) error
	Deregister(service Service) error
	//返回服务
	List(name string) []Service
}
