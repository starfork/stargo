# 04-trace

Demonstrates the tracer interface within handler methods.

## Default behavior

stargo ships with `tracer.DefaultTracer` — a noop implementation.
All trace calls are safe to call; they simply do nothing.

To enable real tracing, swap in a Jaeger tracer before `stargo.New`:

```go
import jtracer "github.com/starfork/stargo/tracer/jaeger"
tracer.DefaultTracer = jtracer.InitJaeger("my-service")
```

The tracer is closed automatically in `app.Stop()`.

## Run

```sh
go run . -c config.yaml
```
