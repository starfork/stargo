# 05-broker

Demonstrates NATS broker usage inside a gRPC handler.

## Pattern

- `NewHandler` subscribes to events on startup (e.g., `user.created`)
- Service methods publish events to the broker
- Topics are auto-prefixed with the app name (e.g., `broker-demo.user.created`)

## Prerequisites

- Running NATS server
- YAML config with `broker.host` set

## Run

```sh
go run . -c config.yaml
```
