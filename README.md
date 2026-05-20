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
| [Quick Start](docs/en/quickstart.md) | [快速开始](docs/zh/quickstart.md) |
| [Usage Guide](docs/en/usage.md) | [使用指南](docs/zh/usage.md) |
| [Configuration](docs/en/config.md) | [配置参考](docs/zh/config.md) |
| [Tools & Setup](docs/en/tools.md) | [工具与环境](docs/zh/tools.md) |
| [Architecture](docs/en/architecture.md) | [架构概览](docs/zh/architecture.md) |

## Samples

See [samples/](samples/) directory for runnable examples, each with a proto-defined
service, handler pattern, and README:

- [01-basic](samples/01-basic/) — Minimal gRPC service with handler struct
- [02-logger](samples/02-logger/) — Structured logging in handler methods
- [03-mysql-redis](samples/03-mysql-redis/) — MySQL repo + Redis cache-aside
- [04-trace](samples/04-trace/) — Tracer interface (noop by default)
- [05-broker](samples/05-broker/) — NATS pub/sub from handler methods
- [06-naming](samples/06-naming/) — etcd registry + resolver
- [07-client](samples/07-client/) — gRPC client with service discovery

## Related projects

- [stargo-examples](https://github.com/starfork/stargo-examples) — Example projects
- [go-slice](https://github.com/starfork/go-slice) — Slice utilities
- [go-crypto](https://github.com/starfork/go-crypto) — Encryption utilities
- [gostar](https://github.com/starfork/gostar) — Project generator
