package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/starfork/stargo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters  sync.Map // map[key]*rate.Limiter
	rate      rate.Limit
	burst     int
	cleanUpIn time.Duration
}

func NewRateLimiter(r rate.Limit, burst int, cleanUpIn time.Duration) *RateLimiter {
	krl := &RateLimiter{
		rate:      r,
		burst:     burst,
		cleanUpIn: cleanUpIn,
	}
	return krl
}

// getLimiter 获取或创建一个限流器
func (k *RateLimiter) getLimiter(key string) *rate.Limiter {
	limiter, ok := k.limiters.Load(key)
	if ok {
		return limiter.(*rate.Limiter)
	}
	newLimiter := rate.NewLimiter(k.rate, k.burst)
	k.limiters.Store(key, newLimiter)
	return newLimiter
}

func GetKey(ctx context.Context) (string, error) {
	key := api.MetaFp(ctx)
	if key == "" {
		key = api.MetaIp(ctx)
	}
	return key, nil
}

// UnaryServerInterceptor 根据 IP 限流
func (k *RateLimiter) UnaryServerInterceptor(getKeyFunc ...func(ctx context.Context) (string, error)) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		f := GetKey
		if len(getKeyFunc) > 0 {
			f = getKeyFunc[0]
		}
		key, err := f(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get rate limit key: %v", err)
		}

		limiter := k.getLimiter(key)
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded for %s", key)
		}

		return handler(ctx, req)
	}
}
