# Quick Start

[中文](../zh/quickstart.md)

## Minimal gRPC server

### 1. Define a handler

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

### 2. Wire it up

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

## With stores (config-first)

If the YAML config has `store.mysql` or `store.redis`, they auto-connect.
Blank-import the packages so their `init()` runs:

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

## Minimal config YAML

```yaml
env: dev
timezone: Asia/Shanghai

server:
  addr: ":9090"
```
