# 04-trace

Demonstrates stargo's tracer abstraction.

By default a noop tracer is used. To enable distributed tracing,
swap in a Jaeger (or other) implementation before `stargo.New`.

## Prerequisites

- Running Jaeger agent (for Jaeger tracer)

## Run

```sh
go run main.go -c config.yaml
```
