package custom

type MarshalOptions struct {
	filterPrefix []string //不处理的路径前缀
	filterUrl    []string //不处理的具体路径
}

type MarshalOption func(o *MarshalOptions)

func WithMarhalPrefix(c []string) MarshalOption {
	return func(o *MarshalOptions) {
		o.filterPrefix = c
	}
}

func WithMarhaUrl(c []string) MarshalOption {
	return func(o *MarshalOptions) {
		o.filterUrl = c
	}
}

type ParserOptions struct {
}

type ParserOption func(o *ParserOptions)
