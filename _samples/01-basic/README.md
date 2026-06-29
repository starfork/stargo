# 01-basic

Minimal stargo gRPC service with a handler implementing proto-defined RPCs.

## Pattern

- `handler` struct embeds `pb.UnimplementedSampleServiceServer` (protoc-generated)
- `NewHandler` receives dependencies (here just `logger.Logger`)
- Methods implement the RPC logic; each receives a `context.Context` and a proto request
- `main.go` loads config, creates the app, registers the handler, and runs

## Run

```sh
go run . -c config.yaml
```
