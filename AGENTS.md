# AGENTS.md — stargo

## Build & test

```sh
# Build root module + all sub-modules (replace directives resolve locally)
go build ./...

# Build all samples
go build ./_samples/...

# Run all tests (within each module that has test files)
go test ./...

# Run a single test package (example)
go test ./store/mysql/

# Run a single sample
cd _samples/01-basic && go run . -c config.yaml
```

There is no root Makefile, no golangci-lint config, and no CI build/test/lint pipeline (only a weekly auto-deps workflow).

## Multi-module structure

This is a **Go monorepo with multiple go.mod files**:

- **Root** (`github.com/starfork/stargo`, Go 1.26.4) — framework core (server, config, logger, store/broker/naming/tracer interfaces, interceptors, api, client, cache, queue, util).
- **Nested sub-modules** — plugin implementations (e.g. `store/mysql`, `broker/nats`, `naming/etcd`, `logger/zap`). Each has its own `go.mod` and **must** contain `replace github.com/starfork/stargo => ../..` (adjust `../` depth based on nesting level).
- **Sample modules** in `_samples/` — independent, runnable projects.

**Critical rule**: when updating dependencies, run `go mod tidy` in both the root **and each changed sub-module**. The CI workflow (`auto-update-deps.yml`) shows the pattern: iterate every non-root `go.mod` directory and run `go get -u ./... && go mod tidy` inside each.

**Exception**: `config/etcd/` is a standalone sub-module with **no** `replace` directive for stargo; it does not depend on the root module.

## Plugin registration pattern (critical for adding new components)

The framework uses `init()` + blank import + factory registry:

1. Framework defines `package.Register(name, factory)` and `package.New(name, config)`.
2. Plugins call `Register()` in `init()`, e.g. `func init() { store.Register("mysql", NewMysql) }`.
3. Users activate plugins via `_ "github.com/starfork/stargo/store/mysql"`.
4. At runtime, `app.initConfig()` calls `NewStore(name, cfg)` which looks up the factory by name.

Key registries: `store.Register`, `broker.Register`, `naming.RegisterRegistry`, `naming.RegisterResolver`, `tracer.Register`, `logger.Register`.

When adding a new plugin:
- Create a sub-directory with its own `go.mod` (use Go 1.26.4).
- Add `replace github.com/starfork/stargo => ../..` (adjust `../` depth).
- Call the `Register()` function in `init()` with a factory that returns the interface.
- Wire the factory name into the YAML config (e.g. `store.mysql.name: mysql` maps to `Register("mysql", ...)`).

## Key environment

- `STARGO_LOG_LEVEL` — sets the default logger level at init time (parseable by `logger.GetLevel()`). Default: `InfoLevel`.
- Config file flag: `-c config.yaml` (defaults to `-c ../config/debug.yaml`).

## Common pitfalls

- **Don't run `go get` or `go mod tidy` only at root** — sub-modules have independent `go.mod` files; they will be out of sync.
- **Tests in sub-modules may need infrastructure** — e.g. `broker/nats/nats_test.go` requires a running NATS server, `cache/redis/redis_test.go` requires Redis.
- **`App.Run()` handles signals and exists the process** — it calls `os.Exit()` on SIGTERM/SIGINT/SIGHUP/SIGQUIT.
- **Reflection is only enabled in non-production** (`s.conf.Env != config.ENV_PRODUCTION`). Don't assume it's always on.
- **gRPC version mismatch risk** — sub-modules pin their own gRPC versions (e.g. `naming/etcd` uses `google.golang.org/grpc v1.79.3` while root uses `v1.76.0`). Do not force-align them without checking etcd compatibility.
- **Retracted versions** — root `go.mod` retracts `[v0.0.1, v0.0.8]` and `[v0.1.1, v0.1.9]`.

## CI

```sh
# Local CI simulation (runs both tasks)
# 1. Build + vet root module
go build -race ./... && go vet ./...

# 2. Build + vet all nested sub-modules
for dir in $(find . -name go.mod ! -path './_samples/*' -exec dirname {} \;); do
  [ "$dir" = "." ] && continue
  (cd "$dir" && go build -race ./... && go vet ./...) || echo "FAIL: $dir"
done

# 3. Build + vet all samples
for dir in $(find _samples -name go.mod -exec dirname {} \;); do
  (cd "$dir" && go build -race ./... && go vet ./...) || echo "FAIL: $dir"
done
```

A full CI workflow exists at `.github/workflows/ci-build.yml` covering root, 14 nested modules, and 21 sample modules with `-race` and `go vet`.

## New features (ROADMAP round)

### Graceful shutdown (A)
`app.Stop()` follows: NOT_SERVING health -> drain in-flight -> Deregister -> GracefulStop. `stopStargo()` respects the order. Health server tracks dependency readiness via `SetDependency()` and exposes `readiness` check endpoint.

### Weighted round-robin balancer (A)
`naming/etcd/balancer.go` registers a `weighted_round_robin_xds` balancer. Services registered with Weight/Version metadata stored in etcd endpoint metadata. `Registry.List()` fully implemented.

### Timeout/deadline interceptors (B)
`interceptor/timeout/timeout.go` provides Unary/Stream server interceptors (default timeout) and Unary/Stream client interceptors (deadline validation). Wired into `server.newRpcServer()` (via `Config.DefaultTimeout`, default 60s) and `client.DefaultOptions()`.

### Readiness health (D)
`server/health.go` supports `readiness` service check that aggregates dependency health (store/broker/registry). `app.initStore/initBroker/initRegistry` register dependencies as healthy.

### mTLS + CORS + TLS (I)
- `internal/tls/tls.go`: `NewServerTransportCredentials()` and `NewClientTransportCredentials()` with CA-based client/server verification.
- `server.Config`: `CertFile`, `KeyFile`, `CAFile` fields; wired into `newRpcServer()`.
- `api/config.go`: `CORS`, `CertFile`, `KeyFile` fields.
- `api/cors.go`: `CORSWrapper()` middleware with configurable origins/methods/headers.
- `api/api.go`: `Run()` supports CORS wrapping and `ListenAndServeTLS`.

## Reference

- Architecture docs: `_docs/en/architecture.md` and `_docs/zh/architecture.md`
- Config reference: `_docs/en/config.md` and `_docs/zh/config.md`
- 14 runnable samples: `_samples/01-basic` through `_samples/14-full-stack`
- TODO / roadmap: `_todo/` (P0 = urgent, P1 = high, P2 = medium, P3 = low)
- Dev reports: `report.md`, `report-dev.md`, `report-plugin.md` at root
