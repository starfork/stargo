# stargo

Go microservice framework: gRPC + etcd (registry/resolver) + NATS (broker) + MySQL/GORM (store) + Redis (store/queue/cache) + Jaeger (tracer) + grpc-gateway (HTTP API).

Module: `github.com/starfork/stargo` · Go 1.25.2 · Apache 2.0. README is in Chinese; documentation and examples at [stargo-examples](https://github.com/starfork/stargo-examples).

## Commands

```sh
go build ./...          # build all packages
go test ./...           # test all packages (some require external infra)
go test ./util/...      # test only utility packages (no external deps)
go test -run TestXxx ./pkg  # single test
```

No Makefile, no linter config, no CI workflows. All packages except the root are leaf libraries with no internal dependencies on each other.

## Test quirks

- **`broker/nats`** and **`cache/redis`** tests require running NATS / Redis — they will panic without them. All other test packages (`util/*`, `pm`, `naming/etcd`, `store/mysql`, `api/custom`) pass standalone.
- `go test -short` is not wired; skip infra-dependent tests by targeting specific packages: `go test ./util/... ./pm/... ./api/custom/...`

## Architecture

- **Entrypoint**: `stargo.New("name", config)` → `app.Run(desc, impl)`. Read `app.go` first.
- **Store access**: `app.Store("mysql").(*mysql.Mysql).GetInstance()` → `*gorm.DB`; same for redis → `*redis.Client`.
- **Stores are opt-in**: mysql/redis packages are NOT imported by `app.go`. They register via `init()`. Users must blank-import them in `main.go`: `_ "github.com/starfork/stargo/store/mysql"` or `_ "github.com/starfork/stargo/store/redis"`.
- **Service discovery**: etcd-backed. Registry at `naming/etcd/registry.go`, resolver at `naming/etcd/resolver.go`. Client uses `round_robin` LB with 1GB/4GB message windows.
- **Interceptors**: auth, zap logging, rate-limit, panic recovery, struct validator (Chinese locale, custom `money` validation).
- **gRPC-gateway**: HTTP→gRPC with optional AES-GCM encrypted marshaler (`api/custom/`). CORS middleware auto-forwards headers as `Grpc-Metadata-*`.
- **Config**: YAML via `config.LoadConfig()`. Default timezone `Asia/Shanghai`. `DefaultConfig` has nil Broker/Registry — no auto-connection attempts.
- **Tracer**: `tracer.DefaultTracer` is a no-op by default. Swap in `tracer/jaeger` or another implementation as needed.
- **Signal handling**: `app.Run()` owns signal handling (SIGTERM/SIGINT/SIGHUP/SIGQUIT) and calls `app.Stop()` for proper cleanup (registry deregister, store close, broker unsubscribe). `server.Server` has no signal handling.

## Key conventions

- Functional options pattern everywhere (see `options.go`, `queue/options.go`, `api/config.go`).
- `pm.Pm` (`map[string]any`) used as the generic parameter bundle — has typed getters (`GetString`, `GetInt`, etc.) and URL encoding.
- Retracted versions `[v0.1.1, v0.1.9]` and `[v0.0.1, v0.0.8]` — do not depend on those.
- Store config supports env var override: `MYSQL_USER`, `MYSQL_PASSWD`, `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_NAME`, `REDIS_HOST`, `REDIS_AUTH`, `REDIS_NUM`.
