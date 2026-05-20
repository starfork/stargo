# Usage Guide

[中文](../zh/usage.md)

## Creating the App

```go
conf, _ := config.LoadConfig()
app := stargo.New("service-name", conf)
```

`stargo.New` initializes stores, broker, registry, resolver, and tracer based on the config.

## Running the server

```go
app.Run(desc, impl)
// or register services first:
// pb.RegisterMyServiceServer(app.RpcServer(), &handler{})
// app.Run()
```

`Run` calls `beforeRun` (sets timezone, creates gRPC server, registers reflection in non-production), then registers services and starts gRPC `Serve`. It owns signal handling for SIGTERM/SIGINT/SIGHUP/SIGQUIT and cleans up on shutdown.

## Stores

Stores are lazily initialized on first access. They must be registered via blank import:

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

Access pattern:

```go
db := app.Store("mysql").(*mysql.Mysql).GetInstance()   // *gorm.DB
rdc := app.Store("redis").(*redis.Redis).GetInstance()   // *redis.Client
```

## Interceptors

Built-in interceptors under `interceptor/`:

| Package | Description |
|---------|-------------|
| `auth` | Bearer/basic token auth with per-service override |
| `logger/zap` | Zap-based request logging |
| `ratelimit` | Per-key rate limiting (fingerprint/IP) |
| `recovery` | Panic recovery returning codes.Unknown |
| `validator` | Struct validation with Chinese locale + custom `money` rule |

Pass interceptors via `grpc.ChainUnaryInterceptor` in server config:

```go
conf.Server.UnaryInterceptor = append(conf.Server.UnaryInterceptor,
    recovery.Unary(),
    validator.Unary(),
)
```

## gRPC-gateway

HTTP/JSON API via grpc-gateway in `api/` package:

```go
import "github.com/starfork/stargo/api"

gw := api.NewApi(&api.Config{
    App:      "my-service",
    Port:     ":8080",
    Registry: namingConfig,
})
```

CORS middleware in `api/cors.go` forwards HTTP headers as `Grpc-Metadata-*`. Optional AES-GCM encrypted marshaler in `api/custom/`.

## Queue (delayed tasks)

```go
import (
	"github.com/starfork/stargo/queue"
	"github.com/starfork/stargo/queue/store/redis"
	"github.com/starfork/stargo/queue/task"
)

store := redis.NewRedis(redisConfig)
q := queue.New(store)
q.Register("my_task", func(t *task.Task) error {
    // process task
    return nil
})
```

## Client (service discovery)

```go
import "github.com/starfork/stargo/client"

c := client.New(ctx, resolver, logger)
conn, err := c.NewClient("target-service")
```
