package timeout

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unary returns a unary server interceptor that enforces a default timeout.
// If the incoming request context has no deadline, sets one using defaultTimeout.
// If it already has a deadline, passes through unchanged.
func Unary(defaultTimeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if defaultTimeout > 0 {
			if _, ok := ctx.Deadline(); !ok {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
				defer cancel()
			}
		}
		return handler(ctx, req)
	}
}

// Stream returns a stream server interceptor that enforces a default timeout.
func Stream(defaultTimeout time.Duration) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		if defaultTimeout > 0 {
			if _, ok := ctx.Deadline(); !ok {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
				defer cancel()
				ss = &wrappedStream{ServerStream: ss, ctx: ctx}
			}
		}
		return handler(srv, ss)
	}
}

// UnaryClient returns a unary client interceptor that propagates deadline.
// If the context has a deadline, it is forwarded to the server via gRPC metadata.
func UnaryClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if deadline, ok := ctx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return status.Errorf(codes.DeadlineExceeded, "deadline already exceeded: %v", remaining)
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClient returns a stream client interceptor that propagates deadline.
func StreamClient() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if deadline, ok := ctx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil, status.Errorf(codes.DeadlineExceeded, "deadline already exceeded: %v", remaining)
			}
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
