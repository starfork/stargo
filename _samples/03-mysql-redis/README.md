# 03-mysql-redis

Demonstrates MySQL and Redis stores within a gRPC handler.

## Config-first

If the YAML config has a `store.mysql` section, stargo auto-connects.
`NewHandler` receives the connected `*gorm.DB` (and optionally `*redis.Client`)
and uses them in service methods.

## Pattern

```go
handler struct {
    repo *repo       // holds db + rdb connections
    log  logger.Logger
    pb.UnimplementedSampleServiceServer
}

repo struct {
    db  *gorm.DB
    rdb *redis.Client
}
```

- `repo` abstracts the data layer
- Handler methods call `repo` for queries and caching
- Redis is used for cache-aside (GetUser checks cache first)

## Prerequisites

- Running MySQL
- Running Redis (optional, for caching)

## Run

```sh
go run . -c config.yaml
```
