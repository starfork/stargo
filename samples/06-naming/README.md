# 06-naming

Demonstrates service registration and discovery via etcd.

The service is automatically registered on `app.Run()` and
deregistered on `app.Stop()`.

## Prerequisites

- Running etcd cluster

## Run

```sh
go run main.go -c config.yaml
```
