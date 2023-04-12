package auth

import (
	"context"

	"google.golang.org/grpc"
)

type AuthFunc func(ctx context.Context) (context.Context, error)
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, method string, req interface{}) (context.Context, interface{}, error)
}

func AuthServerUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		//var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			ctx, req, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod, req)
		}
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// func StreamServerInterceptor() grpc.StreamServerInterceptor {
// 	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		ctx := stream.Context()
// 		var err error
// 		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
// 			ctx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		wrapped := grpc_middleware.WrapServerStream(stream)
// 		wrapped.WrappedContext = ctx
// 		return handler(srv, wrapped)
// 	}
// }
