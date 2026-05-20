# 使用指南

[English](../en/usage.md)

## 配置优先原则

stargo 遵循**配置优先**的原则——只要 YAML 配置中存在对应的配置项，组件会在
`stargo.New` 时自动初始化，无需手动装配。

```go
conf, _ := config.LoadConfig()
app := stargo.New("service-name", conf)
// Stores、broker、registry、resolver、tracer 根据配置自动初始化
```

## Handler 模式

业务逻辑实现在 **handler 结构体**中，它：

1. 嵌入 protoc 生成的 `pb.UnimplementedXxxServer`
2. 通过构造函数 `NewHandler` 接收依赖
3. 实现 RPC 方法

```go
type handler struct {
    repo *repo           // 数据访问层（可选）
    log  logger.Logger
    pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
    h := &handler{log: app.Logger()}
    // 自动连接的 store 通过 app.Store() 获取
    if db := app.Store("mysql"); db != nil {
        h.repo = &repo{db: db.Instance().(*gorm.DB)}
    }
    return h
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    h.log.Infof("GetUser: id=%d", req.Id)
    // 使用 h.repo、h.log 等执行业务逻辑
    return &pb.GetUserResponse{Id: req.Id, Name: "Alice"}, nil
}
```

## 运行服务

```go
conf, _ := config.LoadConfig()
app := stargo.New("my-service", conf)
h := NewHandler(app)

pb.RegisterMyServiceServer(app.RpcServer(), h)
app.Run(pb.MyService_ServiceDesc, h)
```

`app.Run` 执行：
1. 设置时区
2. 创建配置了拦截器的 gRPC 服务
3. 非生产环境注册 gRPC reflection
4. 向 etcd 注册服务（如果配置了 registry）
5. 注册你的服务描述符
6. 安装信号处理（SIGTERM/SIGINT/SIGHUP/SIGQUIT）
7. 调用 `grpc.Serve`

## 存储（MySQL、Redis）

Stores 通过 blank import + YAML 配置按需启用：

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

如果 YAML 中存在 `store.mysql`（或 `store.redis`），store 会在启动时自动连接。
在 handler 中访问：

```go
db := app.Store("mysql").Instance().(*gorm.DB)
rdc := app.Store("redis").Instance().(*redis.Client)
```

## 拦截器

内置拦截器位于 `interceptor/` 目录：

| 包 | 说明 |
|----|------|
| `auth` | Bearer/basic token 认证，支持按服务覆盖 |
| `logger/zap` | 基于 Zap 的请求日志 |
| `ratelimit` | 基于 key 的限流（指纹/IP） |
| `recovery` | panic 恢复，返回 codes.Unknown |
| `validator` | 结构体验证，中文语言包 + 自定义 `money` 规则 |

通过 `grpc.ChainUnaryInterceptor` 传入：

```go
conf.Server.UnaryInterceptor = append(conf.Server.UnaryInterceptor,
    recovery.Unary(),
    validator.Unary(),
)
```

## gRPC-gateway

`api/` 包提供了基于 grpc-gateway 的 HTTP/JSON API：

```go
import "github.com/starfork/stargo/api"

gw := api.NewApi(&api.Config{
    App:      "my-service",
    Port:     ":8080",
    Registry: namingConfig,
})
```

CORS 中间件在 `api/cors.go` 中，将 HTTP 头转发为 `Grpc-Metadata-*`。
`api/custom/` 提供了可选的 AES-GCM 加密 marshaler。

## 队列（延迟任务）

```go
import (
	"github.com/starfork/stargo/queue"
	"github.com/starfork/stargo/queue/store/redis"
	"github.com/starfork/stargo/queue/task"
)

store := redis.NewRedis(redisConfig)
q := queue.New(store)
q.Register("my_task", func(t *task.Task) error {
    // 处理任务
    return nil
})
```

## 客户端（服务发现）

```go
conn, err := app.Client().NewClient("target-service")
```
