## Usage

### main.go

```
import (
	"flag"
	"service/app/internal/server"
	pb "your/proto/path/v1"
 
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/store/mysql"
	"github.com/starfork/stargo/store/redis"
	"google.golang.org/grpc/reflection"
)

func main() { 
	cf := flag.String("c", "../../config/debug.yaml", "config file path")
	flag.Parse()
	sc := server.LoadConfig(*cf)
	c := sc.Server
	app := stargo.New(
		stargo.Org("park"),
		stargo.Name("app"),
		stargo.Config(sc.Server), 
		stargo.UnaryInterceptor(your inteceprot1),
		stargo.UnaryInterceptor(your inteceprot2),
        ...
	)

	s := app.Server()
	if c.Environment == "debug" {
		reflection.Register(s)
	}
	pb.RegisterAppServer(s, server.New(app))
	app.Run()
}
```

### server/handler.go

```
type handler struct {
	logger logger.Logger 
	c      cache.Cache  
	pb.UnimplementedAppServer
}

// New handler
func New(app *stargo.App) *handler {
    rdc:=app.Store("redis").(*redis.Redis).GetInstance() 
	h := &handler{
		logger: app.GetLogger(), 
		c:      credis.New(rdc),
	}  
	return h
}
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

#### slice操作相关

https://github.com/starfork/go-slice

#### 加密

https://github.com/starfork/go-crypto 


