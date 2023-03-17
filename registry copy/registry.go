package registry

type Service struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`    // 服务地址
	Version string `json:"version"` // 服务版本
	Weight  int64  `json:"weight"`  // 服务权重
}

type Registry interface {
	Register(service Service) error
	UnRegister() error
}
