# stargo

Go microservice framework built on gRPC, etcd, NATS, MySQL/GORM, Redis, and Jaeger.

**Module**: `github.com/starfork/stargo` · Go 1.25.2 · Apache 2.0

## Features

- gRPC server with graceful lifecycle
- Service discovery via etcd (registry + resolver)
- Message broker via NATS
- gRPC-gateway HTTP/JSON API with optional AES-GCM encryption
- MySQL/GORM and Redis stores (opt-in via blank import)
- Delayed task queue (Redis sorted sets)
- Configurable interceptors: auth, zap logging, rate-limit, panic recovery, struct validator
- Distributed tracing (Jaeger/OpenTracing)
- Rich utility packages (string, geo, data merging, ID generation, etc.)

## Quick start

```go
package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("my-service", conf)
	// pb.RegisterMyServiceServer(app.RpcServer(), &handler{})
	app.Run()
}
```

## Documentation

| English | 中文 |
|---------|------|
| [Quick Start](_docs/en/quickstart.md) | [快速开始](_docs/zh/quickstart.md) |
| [Usage Guide](_docs/en/usage.md) | [使用指南](_docs/zh/usage.md) |
| [Configuration](_docs/en/config.md) | [配置参考](_docs/zh/config.md) |
| [Tools & Setup](_docs/en/tools.md) | [工具与环境](_docs/zh/tools.md) |
| [Architecture](_docs/en/architecture.md) | [架构概览](_docs/zh/architecture.md) |

## Samples

See [_samples/](_samples/) directory for runnable examples, each with a proto-defined
service, handler pattern, and README:

- [01-basic](_samples/01-basic/) — Minimal gRPC service with handler struct
- [02-logger](_samples/02-logger/) — Structured logging in handler methods
- [03-mysql-redis](_samples/03-mysql-redis/) — MySQL repo + Redis cache-aside
- [04-trace](_samples/04-trace/) — Tracer interface (noop by default)
- [05-broker](_samples/05-broker/) — NATS pub/sub from handler methods
- [06-naming](_samples/06-naming/) — etcd registry + resolver
- [07-client](_samples/07-client/) — gRPC client with service discovery

## Related projects

- [stargo-examples](https://github.com/starfork/stargo-examples) — Example projects
- [go-slice](https://github.com/starfork/go-slice) — Slice utilities
- [go-crypto](https://github.com/starfork/go-crypto) — Encryption utilities
- [gostar](https://github.com/starfork/gostar) — Project generator
