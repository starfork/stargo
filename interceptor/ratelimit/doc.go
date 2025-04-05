package ratelimit

/**
 rateLimiter := interceptor.NewKeyedRateLimiter(5, 10, time.Minute) // 每个 IP 每秒 5 个请求，突发 10 个

server := grpc.NewServer(
	grpc.UnaryInterceptor(rateLimiter.UnaryServerInterceptor(KEY_FUNC OR keep EMPTY),
)
*/
