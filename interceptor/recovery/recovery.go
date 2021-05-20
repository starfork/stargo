package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//Unary interceptor
func Unary() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(Interceptor())
}

// Interceptor panic时返回Unknown错误吗
func Interceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return grpc.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
}
