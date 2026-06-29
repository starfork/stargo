# 工具与环境设置

[English](../en/tools.md)

## Go 工具

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/favadi/protoc-go-inject-tag@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/infobloxopen/protoc-gen-gorm@latest
```

## 调试工具

```sh
go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest
```

## Linting

```sh
go install github.com/yoheimuta/protolint/cmd/protolint@latest
```

## protobuf 安装

| 系统 | 命令 |
|------|------|
| macOS | `brew install protobuf` |
| Debian/Ubuntu | `apt-get install protobuf-compiler` |
| Alpine | `apk add protobuf` |
| Arch Linux | `pacman -S protobuf` |
| CentOS/Fedora | `yum install protobuf-compiler` |
| Docker | `docker run cmd.cat/protoc protoc` |

## 项目生成器

```sh
go install github.com/starfork/gostar@latest
```
