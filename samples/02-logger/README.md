# 02-logger

Demonstrates stargo's logger system.

The default logger outputs to stdout with configurable levels.
Set `STARGO_LOG_LEVEL` env var or call `logger.DefaultLogger.Init()`.

## Run

```sh
STARGO_LOG_LEVEL=debug go run main.go -c config.yaml
```
