package naming

type Registry interface {
	Scheme() string

	Register(service Service) error
	Deregister(service Service) error
	//返回服务
	List(name string) []Service
}


// func NewRegistry(conf *config.Registry) Registry {
// 	name := strings.ToLower(conf.Name)
// 	if name == "redis" {
// 		return redis.NewRegistry(conf)
// 	}
// 	if name == "etcd" {
// 		return etcd.NewRegistry(conf)
// 	}
// 	panic("unknow registry")
// }
