# 02-logger

Demonstrates structured logging inside handler methods.

## Log levels

- `Debugf` — verbose debugging
- `Infof` — normal operational messages
- `Warnf` — non-critical issues
- `Errorf` — errors that need attention

Set `STARGO_LOG_LEVEL=debug` to see all levels.

## Pattern

The handler receives `logger.Logger` from `NewHandler`. Each method uses
the logger naturally as part of its business logic.

## Run

```sh
STARGO_LOG_LEVEL=debug go run . -c config.yaml
```
