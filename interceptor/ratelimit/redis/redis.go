package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/interceptor/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedisRateLimiter struct {
	rdc        *redis.Client
	rate       int64
	windowTime int64
	keyPrefix  string
}

func NewRedisRateLimiter(rdc *redis.Client, rate, windowTime int64, keyPrefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		rdc:        rdc,
		rate:       rate,
		windowTime: windowTime,
		keyPrefix:  keyPrefix,
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string) bool {
	script := `
		local key = KEYS[1]
		local windowTime = ARGV[1]
		local count = ARGV[2]
		local now = ARGV[3]
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, 0, tonumber(now) - tonumber(windowTime))
		
		-- Count current entries
		local current = redis.call('ZCARD', key)
		
		if tonumber(current) < tonumber(count) then
			-- Add current request
			redis.call('ZADD', key, now, now .. math.random())
			redis.call('EXPIRE', key, tonumber(windowTime))
			return 1
		end
		
		return 0
	`
	
	result, err := r.rdc.Eval(ctx, script, []string{r.keyPrefix + key}, r.windowTime, r.rate, time.Now().Unix()).Int()
	if err != nil {
		return false
	}
	
	return result == 1
}

func (r *RedisRateLimiter) UnaryServerInterceptor(getKeyFunc ...func(ctx context.Context) (string, error)) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		f := ratelimit.GetKey
		if len(getKeyFunc) > 0 {
			f = getKeyFunc[0]
		}
		key, err := f(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get rate limit key: %v", err)
		}

		if !r.Allow(ctx, key) {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded for %s", key)
		}

		return handler(ctx, req)
	}
}
