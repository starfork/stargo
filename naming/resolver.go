package naming

import "google.golang.org/grpc/resolver"

type Resolver interface {
	resolver.Builder
}