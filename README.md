## Usage

### Minimal gRPC server (no mysql, no redis)

```go
package main

import (
	"flag"
	"log"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "your/proto/path/v1"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := stargo.New("my-service", conf)

	// gRPC reflection in non-production
	// reflection.Register(app.RpcServer())

	pb.RegisterMyServiceServer(app.RpcServer(), &myHandler{})
	app.Run()
}
```

### With store support (opt-in via blank import)

To use MySQL or Redis stores, add blank imports in your main package:

```go
import (
	_ "github.com/starfork/stargo/store/mysql"  // registers mysql store
	_ "github.com/starfork/stargo/store/redis"  // registers redis store
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("my-service", conf)

	// Access store instances (only if configured in YAML)
	// db := app.Store("mysql").(*mysql.Mysql).GetInstance()  // *gorm.DB
	// rdc := app.Store("redis").(*redis.Redis).GetInstance() // *redis.Client

	app.Run()
}
```

### With interceptors

```go
import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/interceptor/validator"
)

app := stargo.New("my-service", conf)
s := app.RpcServer()

// Register gRPC services
// pb.RegisterMyServiceServer(s, myHandler)

app.Run()
```

### Configuration (YAML)

```yaml
env: dev
timezone: Asia/Shanghai

server:
  addr: ":9090"

# store:
#   mysql:
#     host: localhost
#     port: 3306
#     user: root
#     auth: password
#     name: dbname
#   redis:
#     host: localhost
#     port: 6379
#     auth: ""
#     num: 0

# registry:
#   scheme: etcd
#   host: localhost:2379
#   org: my-org
#   ttl: 10

# broker:
#   host: localhost:4222
```

[更多参考 stargo-examples](https://github.com/starfork/stargo-examples)

## 环境、工具

### protobuf 安装

```

#Debian
apt-get install protobuf-compiler
#Ubuntu
apt-get install protobuf-compiler
#Alpine
apk add protobuf
#Arch Linux
pacman -S protobuf
#Kali Linux
apt-get install protobuf-compiler
#CentOS
yum install protobuf-compiler
#Fedora
dnf install protobuf-compiler

#OS X
brew install protobuf

#Raspbian
apt-get install protobuf-compiler

#Docker
docker run cmd.cat/protoc protoc
```

### 环境相关

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/favadi/protoc-go-inject-tag@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/infobloxopen/protoc-gen-gorm@latest

```

### 工具相关

#### grpc-client-cli 命令行调试工具

```
go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest

```

#### vscocde 插件

```
go install github.com/yoheimuta/protolint/cmd/protolint@latest
```

#### gostar 项目生成工具

```
go install  github.com/starfork/gostar@latest
```

#### 相关库

#### slice 操作相关

https://github.com/starfork/go-slice

#### 加密

https://github.com/starfork/go-crypto
