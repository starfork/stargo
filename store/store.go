package store

type Store interface {
	Instance(conf ...*Config) any
	InstanceE(conf ...*Config) (any, error)
	Close() //关闭连接
}

var TIME_LOCATION = "Asia/Shanghai"
var TFORMAT = "2006-01-02T15:04:05+08:00"
var TZ1K = true

var (
	storeFactories = make(map[string]func(*Config) Store)
)

func Register(name string, factory func(*Config) Store) {
	storeFactories[name] = factory
}

func NewStore(name string, conf *Config) Store {
	if f, ok := storeFactories[name]; ok {
		return f(conf)
	}
	return nil
}
