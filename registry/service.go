package registry

type Service struct {
	Org     string `json:"org"`
	Name    string `json:"name"`
	Addr    string `json:"addr"`    // 服务地址
	Version string `json:"version"` // 服务版本
	Weight  int64  `json:"weight"`  // 服务权重
}
