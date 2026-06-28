# stargo 示例代码 / Stargo Samples

每个示例独立可运行，对应框架的特定特性。所有示例共享 `samples/proto/sample/` 中的 protobuf 定义。

Each sample is self-contained and demonstrates specific framework features. All samples share protobuf definitions from `samples/proto/sample/`.

---

## 特性清单 / Feature Matrix

| # | 示例 / Sample | 特性 / Feature | 核心依赖 / Core | 可选依赖 / Optional |
|---|--------------|----------------|-----------------|-------------------|
| 01 | basic | 最小 gRPC 服务 / Minimal gRPC | 无 / None | 无 / None |
| 02 | logger | 结构化日志 / Structured Logging | 无 / None | zap (contrib) |
| 03 | mysql-redis | 存储层 MySQL + Redis / Stores | 无 / None | mysql, redis (contrib) |
| 04 | trace | 分布式追踪 / Distributed Tracing | 无 / None | jaeger (contrib) |
| 05 | broker | 消息队列 NATS / Message Broker | 无 / None | nats (contrib) |
| 06 | naming | 服务发现 etcd / Service Discovery | 无 / None | etcd (contrib) |
| 07 | client | gRPC 客户端发现 / Client Discovery | 无 / None | etcd (contrib) |
| 08 | gateway | HTTP API 网关 / HTTP Gateway | grpc-gateway | custom marshaler (contrib) |
| 09 | interceptor | 拦截器链 / Interceptor Chain | auth, recovery, ratelimit | validator, zap (contrib) |
| 10 | cache | 缓存抽象 / Cache Abstraction | cache 接口 | redis cache (contrib) |
| 11 | queue | 延迟任务队列 / Delayed Task Queue | queue 引擎 | redis queue (contrib) |
| 12 | mysql-extras | MySQL 扩展特性 / MySQL Extras | 无 / None | mysql (contrib) |

---

## 可选驱动配置示例 / Optional Driver Configs

框架采用 contrib 模块拆分，每个可选驱动通过 blank-import 按需加载。

The framework splits optional drivers into contrib/ — each is loaded on demand via blank-import.

### 仅使用 MySQL / MySQL Only

```go
// main.go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    _ "github.com/starfork/stargo/contrib/store/mysql"  // 按需加载 MySQL
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
    _ "github.com/starfork/stargo/contrib/store/redis"  // 按需加载 Redis
    pb "github.com/starfork/stargo/samples/proto/sample"
)
```

### 零存储 / No Store (纯 gRPC)

```go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    // 不 import 任何 contrib/store/xxx
    // 编译产物最小，无 gorm/mysql/redis 依赖
    // No contrib/store import — minimal binary with zero ORM deps
    pb "github.com/starfork/stargo/samples/proto/sample"
)
```

### 任意组合 / Any Combination

```go
import (
    _ "github.com/starfork/stargo/contrib/store/mysql"
    _ "github.com/starfork/stargo/contrib/store/redis"
    _ "github.com/starfork/stargo/contrib/broker/nats"
    _ "github.com/starfork/stargo/contrib/naming/etcd"
    _ "github.com/starfork/stargo/contrib/tracer/jaeger"
    _ "github.com/starfork/stargo/contrib/interceptor/validator"
    _ "github.com/starfork/stargo/contrib/interceptor/logger/zap"
    // 只 import 你需要的组件
    // Only import what you need
)
```

---

## 快速开始 / Quick Start

```sh
# 构建所有示例 / Build all samples
go build ./samples/...

# 运行特定示例 / Run a specific sample
cd samples/01-basic
go run . -c config.yaml

# 查看可用 flags / See available flags
go run . -h
```

每个示例目录的 `config.yaml` 是完整的配置文件，可根据需要删减不使用的配置段。

Each sample's `config.yaml` is a full config file. Remove unused sections for your use case.
