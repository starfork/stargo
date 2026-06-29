# Tools & Environment Setup

[中文](../zh/tools.md)

## Go tools

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/favadi/protoc-go-inject-tag@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/infobloxopen/protoc-gen-gorm@latest
```

## Debugging

```sh
go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest
```

## Linting

```sh
go install github.com/yoheimuta/protolint/cmd/protolint@latest
```

## protobuf installation

| OS | Command |
|----|---------|
| macOS | `brew install protobuf` |
| Debian/Ubuntu | `apt-get install protobuf-compiler` |
| Alpine | `apk add protobuf` |
| Arch Linux | `pacman -S protobuf` |
| CentOS/Fedora | `yum install protobuf-compiler` |
| Docker | `docker run cmd.cat/protoc protoc` |

## Project generator

```sh
go install github.com/starfork/gostar@latest
```
