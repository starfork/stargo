


### 环境、工具

##### protobuf 安装

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

##### 环境相关

``` 
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/favadi/protoc-go-inject-tag@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest 
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install golang.org/x/tools/cmd/stringer@latest 

```

##### 工具相关

```
go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest

```
vscocde
```
go install github.com/yoheimuta/protolint/cmd/protolint@latest
```