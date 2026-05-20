# 06-naming

Demonstrates etcd service registration and discovery.

## Config-first

When YAML includes a `registry` section with `scheme: etcd`, stargo
auto-creates the registry and resolver.

- **Registration**: happens in `app.Run()` — the service is registered
  with etcd. The handler can access its own service info via `app.Service()`.
- **Deregistration**: happens on `app.Stop()` (SIGTERM/SIGINT).
- **Discovery**: use `app.Resolver()` to build targets for downstream calls.

## Prerequisites

- Running etcd cluster
- YAML config with `registry.scheme: etcd`

## Run

```sh
go run . -c config.yaml
```
