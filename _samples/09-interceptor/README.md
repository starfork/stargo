# 09-interceptor — 拦截器链示例

演示 gRPC 拦截器链的注册与使用：panic recovery、认证、限流。

---

## 拦截器分类

### 内置拦截器（根模块内置，直接 import）

以下拦截器位于 `github.com/starfork/stargo/interceptor/*`，随根模块发布，**无需额外 go get**：

| 包 | 说明 |
|---|------|
| `interceptor/auth` | Bearer token 认证 |
| `interceptor/ratelimit` | 令牌桶限流（支持 Redis 分布式） |
| `interceptor/recovery` | Panic 恢复 |
| `interceptor/timeout` | 客户端/服务端超时传播 |
| `interceptor/bulkhead` | 并发隔离（在途请求限制） |
| `interceptor/circuitbreaker` | 三态熔断器（防雪崩） |
| `interceptor/otelserver` | OTel 服务端追踪 |
| `interceptor/otelclient` | OTel 客户端追踪 |

```go
import (
    "github.com/starfork/stargo/interceptor/auth"
    "github.com/starfork/stargo/interceptor/recovery"
)
```

### 可选拦截器（独立子模块，需显式 go get）

以下拦截器拆分到独立 go.mod 子模块，避免将重依赖（validator v10、zap）强制引入所有项目：

| 模块路径 | 说明 |
|---------|------|
| `github.com/starfork/stargo/interceptor/validator` | 请求参数校验（基于 go-playground/validator） |
| `github.com/starfork/stargo/interceptor/zap` | gRPC 调用日志（基于 Uber zap） |

#### 如何使用子模块拦截器

**方式一：go get（推荐，适用于基于 tag 的依赖管理）**

```bash
go get github.com/starfork/stargo/interceptor/validator@v1.1.1
```

```go
import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/interceptor/validator"
)

func main() {
    conf, _ := config.LoadConfig()
    app := stargo.New("demo", conf)

    app.Init(stargo.WithUnaryInterceptor(validator.Unary()))
    // ...
}
```

**方式二：本地 replace（本地开发用）**

```go
// 项目 go.mod
module github.com/my/project

require (
    github.com/starfork/stargo v1.1.1
    github.com/starfork/stargo/interceptor/validator v1.1.1
)

replace github.com/starfork/stargo => ../path/to/local/stargo
replace github.com/starfork/stargo/interceptor/validator => ../path/to/local/stargo/interceptor/validator
```

**方式三：源码替换（go.work，本地多模块开发）**

```go
// go.work
go 1.26.4

use (
    .
    ../stargo
    ../stargo/interceptor/validator
)
```

---

## 拦截器链顺序

```go
app.Init(
    stargo.WithUnaryInterceptor(recovery.Unary()),        // 1. 最外层：panic 恢复
    stargo.WithUnaryInterceptor(auth.UnaryServerInterceptor(fn)), // 2. 认证
    stargo.WithUnaryInterceptor(ratelimit.UnaryServerInterceptor(fn)), // 3. 限流
    // validator.Unary() 放在 handler 之前
)
// order: recovery → auth → ratelimit → validator → handler
```

---

## 配置文件禁用拦截器

```yaml
# config.yaml
server:
  unary_interceptor: []   # 空列表 = 使用 Init() 注册的拦截器
  # 或显式指定，跳过 Init()：
  # unary_interceptor: ["recovery", "auth", "ratelimit"]
```

---

## 运行

```sh
cd _samples/09-interceptor
go run . -c config.yaml
```
