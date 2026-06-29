# stargo

**Go 微服务框架** — 基于 gRPC + etcd + NATS + MySQL/GORM + Redis + Jaeger，配置驱动，插件按需加载。

**Go microservice framework** — gRPC + etcd + NATS + MySQL/GORM + Redis + Jaeger, config-driven, plugins on demand.

> `github.com/starfork/stargo` · Go 1.25+ · Apache 2.0

---

## 为什么选择 stargo？ / Why stargo?

**痛点 / Pain points**

从头搭建一个 Go 微服务，通常需要手动组合 gRPC、服务发现、消息队列、数据库、缓存、追踪、日志、HTTP 网关等十余个组件，再处理好优雅退出、健康检查、配置管理、拦截器链等基础设施。这往往需要数千行胶水代码，且容易出错。

Building a Go microservice from scratch means manually wiring together gRPC, service discovery, message broker, databases, cache, tracing, logging, HTTP gateway and more — plus infrastructure like graceful shutdown, health checks, config management, and interceptor chains. This takes thousands of lines of boilerplate.

**解法 / Solution**

stargo 将上述所有组件统一为 **配置驱动 + 插件注册 + 按需加载** 的框架。
一行 `stargo.New(name, config)` 自动初始化日志、存储、消息、注册、追踪；一行 `app.Run(desc, impl)` 启动服务并处理信号和优雅退出。

stargo unifies all these components into a **config-driven + plugin registry + on-demand loading** framework. One `stargo.New(name, config)` auto-initializes logging, stores, broker, registry, and tracing; one `app.Run(desc, impl)` starts serving with signal handling and graceful shutdown.

**效率提升 / Productivity gains**

| 维度 | 传统方式 | stargo |
|------|---------|--------|
| 初始化代码量 | ~500-2000 行胶水代码 | ~10 行 |
| 组件切换 | 重写连接逻辑 | 改 YAML + 换 blank import |
| 依赖膨胀 | 全部依赖打进一个 binary | 各插件独立 go.mod，按需编译 |
| 生产就绪 | 需自行实现健康检查/指标/优雅退出 | 内置开箱即用 |
| 团队规范 | 各项目写法不一 | 统一模板，14 个可运行示例 |

---

## 核心特性 / Core Features

### 配置驱动 / Config-Driven

一个 YAML 控制所有组件。无配置的组件不启动，不会报错。

One YAML controls all components. Unconfigured components stay inert — no crash, no error.

```yaml
env: dev
server:
  addr: ":50051"
store:
  mysql:
    host: "localhost"
    port: "3306"
    user: "root"
    auth: "password"
    name: "mydb"
  redis:
    host: "localhost:6379"
broker:
  name: nats
  host: "nats://localhost:4222"
registry:
  scheme: etcd
  host: "localhost:2379"
log:
  driver: zap
  level: debug
api:
  port: ":8080"
```

### 插件按需加载 / Plugins on Demand

每个可选组件（etcd, NATS, MySQL, Redis, Jaeger, zap...）都是独立的 Go module（自有 go.mod），通过 **blank import** 激活。不 import 就不编译进二进制。

Every optional component (etcd, NATS, MySQL, Redis, Jaeger, zap...) is a standalone Go module with its own go.mod, activated via **blank import**. No import = not compiled into your binary.

```go
import (
    // 只需引入你需要的 / Only import what you need
    _ "github.com/starfork/stargo/store/mysql"       // GORM MySQL
    _ "github.com/starfork/stargo/store/redis"       // go-redis
    _ "github.com/starfork/stargo/broker/nats"       // NATS pub/sub
    _ "github.com/starfork/stargo/naming/etcd"       // 服务发现 / service discovery
    _ "github.com/starfork/stargo/tracer/jaeger"     // 分布式追踪 / tracing
    _ "github.com/starfork/stargo/logger/zap"        // Uber zap 日志
    _ "github.com/starfork/stargo/interceptor/validator" // 请求校验
)
```

### 注册表模式 / Registry Pattern

框架定义接口，插件通过 `init()` 自动注册。用户只需要 blank import，无需手动依赖注入。

Framework defines interfaces; plugins self-register via `init()`. Users just blank-import — no manual DI wiring.

```
broker.Register("nats", factory)  →  broker.NewBroker("nats", config)
naming.RegisterRegistry("etcd", factory)  →  naming.NewRegistry("etcd", config)
store.Register("mysql", factory)  →  store.NewStore("mysql", config)
```

### 完整功能栈 / Full Feature Stack

| 模块 / Module | 功能 / Function |
|---------------|-----------------|
| **gRPC Server** | 优雅启动/退出、健康检查、反射、Prometheus 指标、拦截器链 |
| **服务发现 / Discovery** | etcd 注册 + 解析器，自动 lease 续约 |
| **消息代理 / Broker** | NATS JetStream pub/sub，topic 按 app name 前缀 |
| **HTTP 网关 / Gateway** | grpc-gateway JSON API，CORS，支持 AES-GCM 加密 |
| **存储 / Store** | MySQL/GORM、Redis、Postgres — 懒连接、env 覆盖、连接池 |
| **缓存 / Cache** | Redis 缓存、内置文件缓存，支持 Incr/Decr |
| **延迟队列 / Queue** | Redis sorted set 延迟任务，可配重试、并发度 |
| **追踪 / Tracing** | Jaeger/OpenTracing，默认 noop tracer |
| **日志 / Logger** | 默认 console / Go slog / Uber zap，可运行时切换 |
| **拦截器 / Interceptor** | auth / recovery / rate-limit (内置)，validator / zap-logger (可选) |
| **客户端 / Client** | gRPC 客户端，集成 etcd 服务发现 |
| **工具集 / Utils** | `pm.Pm` 泛型参数、字符串工具、数据合并、HTTP 请求辅助 |

---

## 快速开始 / Quick Start

### 最小服务 / Minimal Service

```go
package main

import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    pb "path/to/your/proto"
)

func main() {
    conf, _ := config.LoadConfig()              // 读取 YAML 配置
    app := stargo.New("my-service", conf)        // 初始化所有已配置组件
    h := NewHandler(app.Logger())                // 注入依赖
    app.Run(&pb.MyService_ServiceDesc, h)        // 注册服务并启动（含信号处理）
}
```

```yaml
# config.yaml — 最小配置
env: dev
server:
  addr: ":50051"
```

```sh
go run . -c config.yaml
```

### 获取组件 / Accessing Components

```go
// MySQL
db := app.Store("mysql").(*mysql.Mysql).GetInstance() // *gorm.DB

// Redis
rdb := app.Store("redis").(*redis.Redis).GetInstance()   // *redis.Client

// 消息 / Broker
app.Broker().Publish("topic", msg)

// 服务发现客户端 / Discovery client
conn, _ := app.Client().NewClient("other-service")

// 日志 / Logger
app.Logger().Infof("hello from %s", app.Config().Env)

// 追踪 / Tracer (lifecycle managed by app, cast to use opentracing API)
// tracer := app.Tracer().(*jaeger.JaegerTracer)
```

---

## 运行示例 / Samples

14 个独立可运行示例，覆盖所有特性，位于 [_samples/](_samples/)。

14 self-contained, runnable samples covering all features, under [_samples/](_samples/).

| # | 示例 / Sample | 说明 / Description | 可选依赖 / Optional |
|---|-------|------|----------|
| 01 | [basic](_samples/01-basic/) | 最小 gRPC 服务 / Minimal gRPC service | 无 / None |
| 02 | [logger](_samples/02-logger/) | 结构化日志 / Structured logging in handlers | 无 / None |
| 03 | [mysql-redis](_samples/03-mysql-redis/) | MySQL 仓储 + Redis 缓存 / MySQL repo + Redis cache-aside | store/mysql, store/redis |
| 04 | [trace](_samples/04-trace/) | 分布式追踪 / Distributed tracing | tracer/jaeger |
| 05 | [broker](_samples/05-broker/) | NATS 消息发布/订阅 / NATS pub/sub | broker/nats |
| 06 | [naming](_samples/06-naming/) | etcd 服务注册与发现 / etcd registry + resolver | naming/etcd |
| 07 | [client](_samples/07-client/) | gRPC 客户端服务发现 / Client with service discovery | naming/etcd |
| 08 | [gateway](_samples/08-gateway/) | HTTP API 网关 + 加密 / HTTP Gateway with encryption | api/encrypt |
| 09 | [interceptor](_samples/09-interceptor/) | 拦截器链 / Interceptor chain: auth, recovery, ratelimit, validator, zap | interceptor/validator, interceptor/zap |
| 10 | [cache](_samples/10-cache/) | 缓存抽象层 / Cache abstraction | cache/redis, store/redis |
| 11 | [queue](_samples/11-queue/) | 延迟任务队列 / Delayed task queue (Redis sorted sets) | queue/store/redis, store/redis |
| 12 | [mysql-extras](_samples/12-mysql-extras/) | MySQL 高级特性 / MySQL extras (plugins, geo, UID) | store/mysql |
| 13 | [logger](_samples/13-logger/) | Logger 驱动切换 / Switch between default/slog/zap | logger/zap |
| 14 | [full-stack](_samples/14-full-stack/) | 全栈微服务 / Full-stack demo: 4 services + Docker Compose + K8s | store/mysql, store/redis, naming/etcd |

构建与运行 / Build and Run:

```sh
# 构建所有示例 / Build all
go build ./_samples/...

# 运行单个示例 / Run one sample
cd _samples/01-basic
go run . -c config.yaml

# 全栈示例 / Full-stack demo
cd _samples/14-full-stack && docker compose up -d
```

---

## 模块架构 / Module Architecture

```
stargo (根模块 / root, ~8 直接依赖 / direct deps)
├── config/          # 配置加载 / Config loading (YAML → struct)
│   └── etcd/        # etcd 配置管理（独立 go.mod）
├── server/          # gRPC 服务器 / Server + health + metrics
├── broker/          # 消息代理接口 / Broker interface
│   └── nats/        # NATS 实现（独立 go.mod）
├── naming/          # 注册/解析接口 / Registry + Resolver interfaces
│   └── etcd/        # etcd 实现（独立 go.mod）
├── store/           # 存储接口 / Store interface
│   ├── mysql/       # GORM MySQL（独立 go.mod）
│   ├── redis/       # go-redis（独立 go.mod）
│   └── postgres/    # GORM Postgres（独立 go.mod）
├── cache/           # 缓存接口 / Cache interface
│   ├── redis/       # Redis 缓存（独立 go.mod）
│   └── filecache/   # 内置文件缓存 / Built-in file cache
├── queue/           # 延迟任务引擎 / Delayed task engine
│   └── store/redis/ # Redis 排序集存储（独立 go.mod）
├── tracer/          # 追踪接口 / Tracer interface (noop default)
│   └── jaeger/      # Jaeger 实现（独立 go.mod）
├── logger/          # 日志接口 / Logger interface
│   ├── slog/        # Go 标准库 slog（根模块内置）
│   └── zap/         # Uber zap（独立 go.mod）
├── api/             # gRPC-gateway HTTP API（根模块内置）
│   ├── encrypt/     # AES-GCM 加密 marshaler（根模块内置）
│   └── ratelimit/   # HTTP 限流（独立 go.mod）
├── interceptor/     # 拦截器
│   ├── auth/        # 鉴权（根模块内置）
│   ├── recovery/    # panic 恢复（根模块内置）
│   ├── ratelimit/   # 限流（根模块内置）
│   ├── validator/   # 结构体校验（独立 go.mod）
│   └── zap/         # Zap 日志拦截器（独立 go.mod）
├── client/          # gRPC 客户端（集成服务发现）
└── util/            # 工具集 / Utilities (pm, ustring, request...)
```

每个子模块通过 `replace github.com/starfork/stargo => ../../` 指向根模块，确保依赖版本一致。

Each sub-module uses `replace github.com/starfork/stargo => ../../` to reference the root module, keeping dependency versions consistent.

---

## 配套资源 / Resources

| English | 中文 |
|---------|------|
| [Quick Start](_docs/en/quickstart.md) | [快速开始](_docs/zh/quickstart.md) |
| [Usage Guide](_docs/en/usage.md) | [使用指南](_docs/zh/usage.md) |
| [Configuration](_docs/en/config.md) | [配置参考](_docs/zh/config.md) |
| [Tools & Setup](_docs/en/tools.md) | [工具与环境](_docs/zh/tools.md) |
| [Architecture](_docs/en/architecture.md) | [架构概览](_docs/zh/architecture.md) |

---

## 关联项目 / Related Projects

- [stargo-examples](https://github.com/starfork/stargo-examples) — 示例项目集合 / Example projects
- [go-slice](https://github.com/starfork/go-slice) — 切片工具库 / Slice utilities
- [go-crypto](https://github.com/starfork/go-crypto) — 加密工具库 / Encryption utilities
- [gostar](https://github.com/starfork/gostar) — 项目生成器 / Project generator
