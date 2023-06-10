package config

// Options 参数
type Options struct {
	Store StoreInterface
	Table string
}

// Option Option
type Option func(o *Options)

// Name set server name
func Store(name StoreInterface) Option {
	return func(o *Options) {
		o.Store = name
	}
}
func Table(name string) Option {
	return func(o *Options) {
		o.Table = name
	}
}

func DefaultOptions() Options {
	o := Options{
		//Store: ,
		Table: "config",
	}
	return o
}
