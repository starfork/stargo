# 快速开始

[English](../en/quickstart.md)

## 最小化 gRPC 服务

```go
package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := stargo.New("my-service", conf)
	// 注册你的gRPC服务:
	// pb.RegisterMyServiceServer(app.RpcServer(), &handler{})
	app.Run()
}
```

## 使用拦截器

```go
import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/interceptor/recovery"
	"github.com/starfork/stargo/interceptor/validator"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("my-service", conf)

	s := app.RpcServer()
	// 注册gRPC服务, 然后:
	app.Run()
}
```

## 添加存储 (MySQL, Redis)

存储是 opt-in 的。在 `main.go` 中 blank-import 需要的存储包：

```go
import (
	_ "github.com/starfork/stargo/store/mysql"
	_ "github.com/starfork/stargo/store/redis"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("my-service", conf)

	// 如果YAML中配置了，可以访问存储:
	// gormDB := app.Store("mysql").(*mysql.Mysql).GetInstance()
	// rdc    := app.Store("redis").(*redis.Redis).GetInstance()

	app.Run()
}
```

## 最小配置 YAML

```yaml
env: dev
timezone: Asia/Shanghai

server:
  addr: ":9090"
```
