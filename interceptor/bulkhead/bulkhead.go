package bulkhead

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Bulkhead struct {
	sem     chan struct{}
	max     int64
	timeout time.Duration
}

func New(opts ...Option) *Bulkhead {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(cfg)
	}
	return &Bulkhead{
		sem:     make(chan struct{}, cfg.MaxConcurrent),
		max:     cfg.MaxConcurrent,
		timeout: cfg.WaitTimeout,
	}
}

func (b *Bulkhead) Allow(ctx context.Context) error {
	select {
	case b.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if b.timeout > 0 {
		timer := time.NewTimer(b.timeout)
		defer timer.Stop()
		select {
		case b.sem <- struct{}{}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return status.Errorf(codes.ResourceExhausted, "bulkhead: too many concurrent requests")
		}
	}
	return status.Errorf(codes.ResourceExhausted, "bulkhead: too many concurrent requests")
}

func (b *Bulkhead) Release() {
	select {
	case <-b.sem:
	default:
	}
}

func (b *Bulkhead) InFlight() int {
	return len(b.sem)
}

func (b *Bulkhead) Max() int64 {
	return b.max
}

func UnaryClientInterceptor(bh *Bulkhead) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := bh.Allow(ctx); err != nil {
			return err
		}
		defer bh.Release()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(bh *Bulkhead) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if err := bh.Allow(ctx); err != nil {
			return nil, err
		}
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			bh.Release()
		}
		return &bulkheadClientStream{cs, bh}, err
	}
}

type bulkheadClientStream struct {
	grpc.ClientStream
	bh *Bulkhead
}

func (s *bulkheadClientStream) CloseSend() error {
	err := s.ClientStream.CloseSend()
	s.bh.Release()
	return err
}

func UnaryServerInterceptor(bh *Bulkhead) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := bh.Allow(ctx); err != nil {
			return nil, err
		}
		defer bh.Release()
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(bh *Bulkhead) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := bh.Allow(ss.Context()); err != nil {
			return err
		}
		defer bh.Release()
		return handler(srv, ss)
	}
}
