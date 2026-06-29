package encrypt

type MarshalOptions struct {
	filterPrefix []string
	filterUrl    []string
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
