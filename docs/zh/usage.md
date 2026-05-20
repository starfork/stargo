# 使用指南

[English](../en/usage.md)

## 创建应用

```go
conf, _ := config.LoadConfig()
app := stargo.New("service-name", conf)
```

`stargo.New` 根据配置初始化 stores、broker、registry、resolver 和 tracer。

## 运行服务

```go
app.Run(desc, impl)
// 或先注册服务:
// pb.RegisterMyServiceServer(app.RpcServer(), &handler{})
// app.Run()
```

`Run` 会调用 `beforeRun`（设置时区、创建 gRPC 服务、非生产环境注册 reflection），然后注册服务并启动 gRPC `Serve`。它负责信号处理（SIGTERM/SIGINT/SIGHUP/SIGQUIT）并在关闭时清理资源。

## 存储

Stores 在首次访问时懒初始化。必须通过 blank import 注册：

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

访问方式：

```go
db := app.Store("mysql").(*mysql.Mysql).GetInstance()   // *gorm.DB
rdc := app.Store("redis").(*redis.Redis).GetInstance()   // *redis.Client
```

## 拦截器

内置拦截器位于 `interceptor/` 目录：

| 包 | 说明 |
|----|------|
| `auth` | Bearer/basic token 认证，支持按服务覆盖 |
| `logger/zap` | 基于 Zap 的请求日志 |
| `ratelimit` | 基于 key 的限流（指纹/IP） |
| `recovery` | panic 恢复，返回 codes.Unknown |
| `validator` | 结构体验证，中文语言包 + 自定义 `money` 规则 |

通过 `grpc.ChainUnaryInterceptor` 传入拦截器：

```go
conf.Server.UnaryInterceptor = append(conf.Server.UnaryInterceptor,
    recovery.Unary(),
    validator.Unary(),
)
```

## gRPC-gateway

`api/` 包提供了基于 grpc-gateway 的 HTTP/JSON API：

```go
import "github.com/starfork/stargo/api"

gw := api.NewApi(&api.Config{
    App:      "my-service",
    Port:     ":8080",
    Registry: namingConfig,
})
```

CORS 中间件在 `api/cors.go` 中，将 HTTP 头转发为 `Grpc-Metadata-*`。`api/custom/` 提供了可选的 AES-GCM 加密 marshaler。

## 队列（延迟任务）

```go
import (
	"github.com/starfork/stargo/queue"
	"github.com/starfork/stargo/queue/store/redis"
	"github.com/starfork/stargo/queue/task"
)

store := redis.NewRedis(redisConfig)
q := queue.New(store)
q.Register("my_task", func(t *task.Task) error {
    // 处理任务
    return nil
})
```

## 客户端（服务发现）

```go
import "github.com/starfork/stargo/client"

c := client.New(ctx, resolver, logger)
conn, err := c.NewClient("target-service")
```
