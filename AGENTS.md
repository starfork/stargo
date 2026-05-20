# stargo

Go microservice framework: gRPC + etcd (registry/resolver) + NATS (broker) + MySQL/GORM (store) + Redis (store/queue/cache) + Jaeger (tracer) + grpc-gateway (HTTP API).

Root module: `github.com/starfork/stargo` · Contrib module: `github.com/starfork/stargo/contrib` · Go 1.25.2 · Apache 2.0. README is in Chinese; documentation and examples at [stargo-examples](https://github.com/starfork/stargo-examples).

## Module split

The root `go.mod` is minimal (~8 direct deps: grpc, grpc-gateway, grpc-middleware, protobuf, yaml, x/time, x/text, x/exp). All optional/pluggable implementations live in `contrib/` with their own `go.mod`:

| Root (core) | Contrib (optional) |
|---|---|
| `broker/broker.go` — interface + registry | `contrib/broker/nats` — NATS implementation |
| `naming/registry.go`, `resolver.go` — interfaces + registries | `contrib/naming/etcd` — etcd implementation |
| `store/store.go` — interface + registry | `contrib/store/mysql`, `contrib/store/redis` |
| `tracer/tracer.go` — interface (noop default) | `contrib/tracer/jaeger` |
| `interceptor/auth`, `recovery`, `ratelimit` | `contrib/interceptor/validator`, `contrib/interceptor/logger/zap` |
| `api/api.go` — gateway runtime | `contrib/api/custom` — encrypted marshaler |
| `cache/cache.go` — interface | `contrib/cache/redis` |
| `queue/store/store.go` — interface | `contrib/queue/store/redis` |
| `util/request` — HTTP helpers | `contrib/util/request/limiter` — uber/ratelimit |

## Registry pattern

Optional implementations self-register via `init()`:

- **Broker**: `broker.Register("nats", factory)` → used by `broker.NewBroker("nats", config)`
- **Registry**: `naming.RegisterRegistry("etcd", factory)` → used by `naming.NewRegistry("etcd", config)`
- **Resolver**: `naming.RegisterResolver("etcd", factory)` → used by `naming.NewResolver("etcd", config)`
- **Store**: `store.Register("mysql", factory)` → used by `store.NewStore("mysql", config)`

Users blank-import the contrib package in `main.go` to trigger registration:

```go
import _ "github.com/starfork/stargo/contrib/broker/nats"
import _ "github.com/starfork/stargo/contrib/naming/etcd"
import _ "github.com/starfork/stargo/contrib/store/mysql"
```

## Commands

```sh
go build ./...          # build root module
go build ./contrib/...  # build contrib module
go test ./util/...      # test utility packages (no external deps)
go test ./pm/...        # test pm package
```

## Test quirks

- **`contrib/broker/nats`** and **`contrib/cache/redis`** tests require running NATS / Redis — they will panic without them. All other test packages pass standalone.
- Skip infra-dependent tests by targeting specific packages: `go test ./util/... ./pm/...`

## Architecture

- **Entrypoint**: `stargo.New("name", config)` → `app.Run(desc, impl)`. Read `app.go` first.
- **Store access**: `app.Store("mysql").(*mysql.Mysql).GetInstance()` → `*gorm.DB`; same for redis → `*redis.Client`.
- **Stores are opt-in**: mysql/redis packages are NOT imported by `app.go`. They register via `init()`. Users must blank-import contrib versions in `main.go`: `_ "github.com/starfork/stargo/contrib/store/mysql"` or `_ "github.com/starfork/stargo/contrib/store/redis"`.
- **Service discovery**: registry + resolver use registry pattern. Users blank-import `contrib/naming/etcd` to register etcd implementations.
- **Broker**: users blank-import `contrib/broker/nats` to register NATS broker.
- **Interceptors**: auth, rate-limit, panic recovery are core (`interceptor/`). Validator and zap logger are in contrib.
- **gRPC-gateway**: HTTP→gRPC with optional AES-GCM encrypted marshaler (`contrib/api/custom/`). CORS middleware auto-forwards headers as `Grpc-Metadata-*`.
- **Config**: YAML via `config.LoadConfig()`. Default timezone `Asia/Shanghai`. `DefaultConfig` has nil Broker/Registry — no auto-connection attempts.
- **Tracer**: `tracer.DefaultTracer` is a no-op by default. Swap in `contrib/tracer/jaeger` or another implementation as needed.
- **Signal handling**: `app.Run()` owns signal handling (SIGTERM/SIGINT/SIGHUP/SIGQUIT) and calls `app.Stop()` for proper cleanup (registry deregister, store close, broker unsubscribe). `server.Server` has no signal handling.

## Key conventions

- Functional options pattern everywhere (see `options.go`, `queue/options.go`, `api/config.go`).
- `pm.Pm` (`map[string]any`) used as the generic parameter bundle — has typed getters (`GetString`, `GetInt`, etc.) and URL encoding.
- Retracted versions `[v0.1.1, v0.1.9]` and `[v0.0.1, v0.0.8]` — do not depend on those.
- Store config supports env var override: `MYSQL_USER`, `MYSQL_PASSWD`, `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_NAME`, `REDIS_HOST`, `REDIS_AUTH`, `REDIS_NUM`.
