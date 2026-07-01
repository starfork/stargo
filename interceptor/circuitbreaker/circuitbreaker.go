package circuitbreaker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type State int32

const (
	StateClosed   State = 0
	StateOpen     State = 1
	StateHalfOpen State = 2
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

type CircuitBreaker struct {
	state    atomic.Int32
	mu       sync.RWMutex
	openedAt time.Time

	failureThreshold float64
	successThreshold int64
	openTimeout      time.Duration
	halfOpenMaxReqs  int64

	window *slidingWindow

	halfOpenCount atomic.Int64

	onStateChange func(from, to State)
}

func New(opts ...Option) *CircuitBreaker {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	cb := &CircuitBreaker{
		failureThreshold: cfg.FailureThreshold,
		successThreshold: cfg.SuccessThreshold,
		openTimeout:      cfg.OpenTimeout,
		halfOpenMaxReqs:  cfg.HalfOpenMaxReqs,
		window:           newSlidingWindow(cfg.WindowSize),
		onStateChange:    cfg.OnStateChange,
	}
	cb.state.Store(int32(StateClosed))
	return cb
}

func (cb *CircuitBreaker) Allow() error {
	st := cb.State()
	switch st {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(cb.openedAt) > cb.openTimeout {
			cb.toHalfOpen()
			return cb.Allow()
		}
		return status.Errorf(codes.Unavailable, "circuit breaker is OPEN")
	case StateHalfOpen:
		if cb.halfOpenCount.Add(1) > cb.halfOpenMaxReqs {
			cb.halfOpenCount.Add(-1)
			return status.Errorf(codes.Unavailable, "circuit breaker: too many half-open probes")
		}
		return nil
	default:
		return nil
	}
}

func (cb *CircuitBreaker) MarkSuccess() {
	cb.window.addSuccess()

	st := cb.State()
	if st == StateHalfOpen {
		cb.halfOpenCount.Add(-1)
		if cb.window.successCount() >= cb.successThreshold {
			cb.toClosed()
		}
	}
}

func (cb *CircuitBreaker) MarkFailure() {
	cb.window.addFailure()

	st := cb.State()
	switch st {
	case StateClosed:
		if cb.window.failureRate() >= cb.failureThreshold && cb.window.totalRequests() >= 10 {
			cb.toOpen()
		}
	case StateHalfOpen:
		cb.halfOpenCount.Add(-1)
		cb.toOpen()
	}
}

func (cb *CircuitBreaker) State() State {
	return State(cb.state.Load())
}

func (cb *CircuitBreaker) toOpen() {
	old := State(cb.state.Swap(int32(StateOpen)))
	if old == StateOpen {
		return
	}
	cb.mu.Lock()
	cb.openedAt = time.Now()
	cb.mu.Unlock()
	cb.halfOpenCount.Store(0)
	if cb.onStateChange != nil {
		cb.onStateChange(old, StateOpen)
	}
}

func (cb *CircuitBreaker) toHalfOpen() {
	old := State(cb.state.Swap(int32(StateHalfOpen)))
	if old == StateHalfOpen {
		return
	}
	cb.halfOpenCount.Store(0)
	if cb.onStateChange != nil {
		cb.onStateChange(old, StateHalfOpen)
	}
}

func (cb *CircuitBreaker) toClosed() {
	old := State(cb.state.Swap(int32(StateClosed)))
	if old == StateClosed {
		return
	}
	cb.window.reset()
	cb.halfOpenCount.Store(0)
	if cb.onStateChange != nil {
		cb.onStateChange(old, StateClosed)
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.state.Store(int32(StateClosed))
	cb.window.reset()
	cb.halfOpenCount.Store(0)
}

func isServerError(err error) bool {
	if err == nil {
		return false
	}
	code := status.Code(err)
	switch code {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Internal,
		codes.Unknown:
		return true
	default:
		return false
	}
}

func UnaryClientInterceptor(cb *CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := cb.Allow(); err != nil {
			return err
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if isServerError(err) {
			cb.MarkFailure()
		} else {
			cb.MarkSuccess()
		}
		return err
	}
}

func StreamClientInterceptor(cb *CircuitBreaker) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if err := cb.Allow(); err != nil {
			return nil, err
		}
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if isServerError(err) {
			cb.MarkFailure()
		} else {
			cb.MarkSuccess()
		}
		return cs, err
	}
}
