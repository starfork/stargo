package main

import (
	"context"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/interceptor/auth"
	"github.com/starfork/stargo/interceptor/ratelimit"
	"github.com/starfork/stargo/interceptor/recovery"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("interceptor-demo", conf)

	// 注册拦截器 / Register interceptors
	// 执行顺序: recovery → auth → ratelimit → handler
	// Order: recovery → auth → ratelimit → handler

	// 1. panic 恢复 / Panic recovery — 确保 panic 不会杀死进程
	app.Init(stargo.WithUnaryInterceptor(recovery.Unary()))

	// 2. Bearer token 认证 / Bearer token authentication
	// 自定义验证函数：检查 token 是否为 "valid-token"
	// Custom auth function: check if token equals "valid-token"
	app.Init(stargo.WithUnaryInterceptor(
		auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
			token, err := auth.AuthFromMD(ctx, "bearer")
			if err != nil {
				return ctx, status.Error(codes.Unauthenticated, "missing bearer token")
			}
			if token != "valid-token" {
				return ctx, status.Error(codes.PermissionDenied, "invalid token")
			}
			return ctx, nil
		}),
	))

	// 3. 速率限制 / Rate limiting — 每秒最多 10 个请求
	rl := ratelimit.NewRateLimiter(10, 10, 0)
	app.Init(stargo.WithUnaryInterceptor(
		rl.UnaryServerInterceptor(func(ctx context.Context) (string, error) {
			// 按客户端 IP 限流 / Rate limit by client IP
			return "", nil
		}),
	))

	h := NewHandler(app.Logger())
	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
