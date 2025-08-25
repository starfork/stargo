package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unary interceptor
func Unary() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(Interceptor())
}

// Interceptor panic时返回Unknown错误吗
func Interceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
}
