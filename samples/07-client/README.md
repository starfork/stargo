# 07-client

Demonstrates gRPC client with service discovery.

`app.Client()` returns a client that connects to services
discovered via the configured resolver (e.g. etcd).

## Prerequisites

- Running etcd cluster
- A target service registered in etcd

## Run

```sh
go run main.go -c config.yaml
```
