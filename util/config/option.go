package config

// Options 参数
type Options struct {
	Store StoreInterface
}

// Option Option
type Option func(o *Options)

// Name set server name
func Store(name StoreInterface) Option {
	return func(o *Options) {
		o.Store = name
	}
}

func DefaultOptions() Options {
	o := Options{
		//Store: ,
	}
	return o
}
