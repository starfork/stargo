//实现接口
//"google.golang.org/grpc/resolver"

package resolver

import "google.golang.org/grpc/resolver"

const (
	poolScheme         = "example"
	exampleServiceName = "resolver.example.grpc.io"

	backendAddr = "localhost:50051"
)

type poolResolverBuilder struct{}

func NewPoolBuilder() (resolver.Builder, error) {
	return &poolResolverBuilder{}, nil
}

func (*poolResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &poolResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: {backendAddr},
		},
	}
	r.start()
	return r, nil
}
func (*poolResolverBuilder) Scheme() string { return poolScheme }

// poolResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type poolResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *poolResolver) start() {
	addrStrs := r.addrsStore[r.target.URL.Path]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*poolResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*poolResolver) Close()                                  {}
