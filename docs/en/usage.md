# Usage Guide

[中文](../zh/usage.md)

## Config-first principle

stargo follows a **config-first** approach — if a component is present in the YAML
configuration, it is auto-initialized when `stargo.New` is called. No manual wiring
is needed for most components.

```go
conf, _ := config.LoadConfig()
app := stargo.New("service-name", conf)
// Stores, broker, registry, resolver, tracer are auto-initialized
// based on conf contents.
```

## Handler pattern

Service logic lives in a **handler struct** that:

1. Embeds the protoc-generated `pb.UnimplementedXxxServer`
2. Receives dependencies through a constructor (`NewHandler`)
3. Implements the RPC methods

```go
type handler struct {
    repo *repo           // data access layer (optional)
    log  logger.Logger
    pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
    h := &handler{log: app.Logger()}
    // Auto-connected stores are available via app.Store().
    if db := app.Store("mysql"); db != nil {
        h.repo = &repo{db: db.Instance().(*gorm.DB)}
    }
    return h
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    h.log.Infof("GetUser: id=%d", req.Id)
    // Business logic using h.repo, h.log, etc.
    return &pb.GetUserResponse{Id: req.Id, Name: "Alice"}, nil
}
```

## Running the server

```go
conf, _ := config.LoadConfig()
app := stargo.New("my-service", conf)
h := NewHandler(app)

// Register your gRPC service and start serving.
pb.RegisterMyServiceServer(app.RpcServer(), h)
app.Run(pb.MyService_ServiceDesc, h)
```

`app.Run`:
1. Sets the timezone
2. Creates the gRPC server with configured interceptors
3. Registers gRPC reflection in non-production environments
4. Registers the service with etcd (if registry is configured)
5. Registers your service descriptor
6. Installs signal handlers (SIGTERM/SIGINT/SIGHUP/SIGQUIT)
7. Calls `grpc.Serve`

## Stores (MySQL, Redis)

Stores are **opt-in** via blank import + YAML config.

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

If `store.mysql` (or `store.redis`) exists in YAML, the store is auto-connected
at startup. Access it in your handler:

```go
db := app.Store("mysql").Instance().(*gorm.DB)   // *gorm.DB
rdc := app.Store("redis").Instance().(*redis.Client) // *redis.Client
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

CORS middleware in `api/cors.go` forwards HTTP headers as `Grpc-Metadata-*`.
Optional AES-GCM encrypted marshaler in `api/custom/`.

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
conn, err := app.Client().NewClient("target-service")
```
