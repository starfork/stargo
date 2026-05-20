# 03-mysql-redis

Demonstrates MySQL and Redis store usage.

Stores are opt-in — you must blank-import the store packages:

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

## Prerequisites

- Running MySQL instance
- Running Redis instance

## Run

```sh
go run main.go -c config.yaml
```
