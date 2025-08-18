package path

type Options struct {
	base     int //转换机制，一旦程序运行，则不可修改
	maxLevel int //最大层级
}

type Option func(o *Options)

func WithBase(c int) Option {
	return func(o *Options) {
		o.base = c
	}
}

func WithMaxLevel(c int) Option {
	return func(o *Options) {
		o.maxLevel = c
	}
}

func DefaultOptions() *Options {
	o := &Options{
		base:     36,
		maxLevel: 10,
	}
	return o
}
