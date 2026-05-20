# 07-client

Demonstrates gRPC client with service discovery within a handler.

## Pattern

- `app.Client()` returns a client connected via the configured resolver
- Handler methods call `cli.NewClient("target-service")` with the
  service name to discover via etcd
- The returned `*grpc.ClientConn` is used to create a downstream gRPC stub

## Prerequisites

- Running etcd cluster with registering services
- YAML config with `registry.scheme: etcd`

## Run

```sh
go run . -c config.yaml
```
