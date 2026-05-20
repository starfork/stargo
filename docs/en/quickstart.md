# Quick Start

[中文](../zh/quickstart.md)

## Minimal gRPC server

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
	// Register your gRPC service:
	// pb.RegisterMyServiceServer(app.RpcServer(), &handler{})
	app.Run()
}
```

## With interceptors

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
	// Register gRPC services, then:
	app.Run()
}
```

## Adding stores (MySQL, Redis)

Stores are opt-in. Blank-import the store package in your `main.go`:

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

	// Access stores if configured in YAML:
	// gormDB := app.Store("mysql").(*mysql.Mysql).GetInstance()
	// rdc    := app.Store("redis").(*redis.Redis).GetInstance()

	app.Run()
}
```

## Minimal config YAML

```yaml
env: dev
timezone: Asia/Shanghai

server:
  addr: ":9090"
```
