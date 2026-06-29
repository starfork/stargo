# stargo 示例代码 / Stargo Samples

每个示例独立可运行，对应框架的特定特性。所有示例共享 `samples/proto/sample/` 中的 protobuf 定义。

Each sample is self-contained and demonstrates specific framework features. All samples share protobuf definitions from `samples/proto/sample/`.

---

## 特性清单 / Feature Matrix

| # | 示例 / Sample | 特性 / Feature | 核心依赖 / Core | 可选依赖 / Optional |
|---|--------------|----------------|-----------------|-------------------|
| 01 | basic | 最小 gRPC 服务 / Minimal gRPC | 无 / None | 无 / None |
| 02 | logger | 结构化日志 / Structured Logging | 无 / None | 无 / None |
| 03 | mysql-redis | 存储层 MySQL + Redis / Stores | 无 / None | `store/mysql`, `store/redis` |
| 04 | trace | 分布式追踪 / Distributed Tracing | 无 / None | `tracer/jaeger` |
| 05 | broker | 消息队列 NATS / Message Broker | 无 / None | `broker/nats` |
| 06 | naming | 服务发现 etcd / Service Discovery | 无 / None | `naming/etcd` |
| 07 | client | gRPC 客户端发现 / Client Discovery | 无 / None | `naming/etcd` |
| 08 | gateway | HTTP API 网关 / HTTP Gateway | grpc-gateway | `api/encrypt` |
| 09 | interceptor | 拦截器链 / Interceptor Chain | auth, recovery, ratelimit | `interceptor/validator`, `interceptor/zap` |
| 10 | cache | 缓存抽象 / Cache Abstraction | cache 接口 | `cache/redis`, `store/redis` |
| 11 | queue | 延迟任务队列 / Delayed Task Queue | queue 引擎 | `queue/store/redis`, `store/redis` |
| 12 | mysql-extras | MySQL 扩展特性 / MySQL Extras | 无 / None | `store/mysql` |
| 13 | logger | Logger 驱动切换 (default/slog/zap) | logger 接口 | `logger/zap` |
| 14 | full-stack | 全栈微服务 (4服务+Docker+K8s) | 综合 | `store/mysql`, `store/redis`, `naming/etcd`, `logger/zap` |

---

## 可选驱动配置示例 / Optional Driver Configs

框架采用插件子模块拆分，每个可选驱动通过 blank-import 按需加载。未 import 的插件不参与编译。

The framework splits optional drivers into plugin sub-modules — each is loaded on demand via blank-import.

### 仅使用 MySQL / MySQL Only

```go
// main.go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    _ "github.com/starfork/stargo/store/mysql"  // 按需加载 MySQL
    pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
    conf, _ := config.LoadConfig()
    app := stargo.New("mysql-only", conf)
    // 只引入 MySQL, Redis 不会参与编译
    // Only MySQL is compiled in; Redis is excluded
}
```

### 仅使用 Redis / Redis Only

```go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    _ "github.com/starfork/stargo/store/redis"  // 按需加载 Redis
    pb "github.com/starfork/stargo/samples/proto/sample"
)
```

### 零存储 / No Store (纯 gRPC)

```go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    // 不 import 任何 store/xxx
    // 编译产物最小，无 gorm/mysql/redis 依赖
    // No store import — minimal binary with zero ORM deps
    pb "github.com/starfork/stargo/samples/proto/sample"
)
```

### 任意组合 / Any Combination

```go
import (
    _ "github.com/starfork/stargo/store/mysql"
    _ "github.com/starfork/stargo/store/redis"
    _ "github.com/starfork/stargo/broker/nats"
    _ "github.com/starfork/stargo/naming/etcd"
    _ "github.com/starfork/stargo/tracer/jaeger"
    _ "github.com/starfork/stargo/logger/zap"
    _ "github.com/starfork/stargo/interceptor/validator"
    _ "github.com/starfork/stargo/interceptor/zap"
    // 只 import 你需要的组件
    // Only import what you need
)
```

---

## Logger 驱动 / Logger Drivers

```yaml
# config.yaml
log:
  driver: ""      # ""=默认console, "slog"=Go标准库, "zap"=Uber zap
  level: debug
```

| driver | import | 说明 / Note |
|---|---|---|
| `""` | 无 / None | 默认 console 日志 |
| `slog` | 无 / None | Go stdlib `log/slog`，根模块内置 |
| `zap` | `_ "github.com/starfork/stargo/logger/zap"` | Uber 高性能结构化日志 |

---

## 快速开始 / Quick Start

```sh
# 构建所有示例 / Build all samples
go build ./samples/...

# 运行特定示例 / Run a specific sample
cd samples/01-basic
go run . -c config.yaml

# 运行全栈示例 / Run full-stack demo
cd samples/14-full-stack && docker compose up -d

# 查看可用 flags / See available flags
go run . -h
```

每个示例目录的 `config.yaml` 是完整的配置文件，可根据需要删减不使用的配置段。

Each sample's `config.yaml` is a full config file. Remove unused sections for your use case.
