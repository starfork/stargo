# 快速开始

[English](../en/quickstart.md)

## 最小化 gRPC 服务

### 1. 定义 handler

```go
type handler struct {
    pb.UnimplementedSampleServiceServer
}

func NewHandler() *handler {
    return &handler{}
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    return &pb.GetUserResponse{Id: req.Id, Name: "Alice"}, nil
}
```

### 2. 装配启动

```go
package main

import (
    "github.com/starfork/stargo"
    "github.com/starfork/stargo/config"
    pb "your/proto/package"
)

func main() {
    conf, _ := config.LoadConfig()
    app := stargo.New("my-service", conf)
    h := NewHandler()

    pb.RegisterSampleServiceServer(app.RpcServer(), h)
    app.Run(pb.SampleService_ServiceDesc, h)
}
```

## 使用存储（配置优先）

如果 YAML 配置中有 `store.mysql` 或 `store.redis`，它们会自动连接。
通过 blank import 让 `init()` 注册：

```go
import (
    _ "github.com/starfork/stargo/store/mysql"
    _ "github.com/starfork/stargo/store/redis"
)

func NewHandler(app *stargo.App) *handler {
    h := &handler{log: app.Logger()}
    if db := app.Store("mysql"); db != nil {
        h.repo = &repo{db: db.Instance().(*gorm.DB)}
    }
    return h
}
```

## 最小配置 YAML

```yaml
env: dev
timezone: Asia/Shanghai

server:
  addr: ":9090"
```
