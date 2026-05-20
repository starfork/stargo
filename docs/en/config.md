# Configuration Reference

[中文](../zh/config.md)

## YAML structure

```yaml
env: dev                   # dev | docker | production
timezone: Asia/Shanghai    # default timezone
timeformat: "2006-01-02T15:04:05+08:00"

server:
  addr: ":9090"
  # UnaryInterceptor:     # configured programmatically
  # StreamInterceptor:    # configured programmatically
  # ServerOpts:           # configured programmatically

api:
  app: my-service
  port: ":8080"
  registry:               # etcd config for gateway
    scheme: etcd
    host: localhost:2379
  enc: false              # enable AES-GCM encryption
  enckey: ""              # encryption key (32 bytes)

store:
  mysql:
    host: localhost
    port: "3306"
    user: root
    auth: password
    name: dbname
    prefix: myapp_
    debug: false
    max_idle: 10
    max_open: 100

  redis:
    host: localhost
    port: "6379"
    auth: ""
    num: 0                # database number
    debug: false

log:
  level: info             # trace | debug | info | warn | error | fatal

broker:
  host: localhost:4222

registry:
  scheme: etcd
  org: my-org
  host: localhost:2379
  auth: ""
  ttl: 10                 # lease TTL in seconds

tracer:
  host: localhost:6831    # Jaeger agent UDP
  name: my-service

filemanager:
  endpoint: s3.amazonaws.com
  access_key: ""
  secret_key: ""
  bucket_name: my-bucket

jwt:
  public_key: ""
  private_key: ""
```

## Environment variable overrides

| Var | Overrides |
|-----|-----------|
| `MYSQL_USER` | store.mysql.user |
| `MYSQL_PASSWD` | store.mysql.auth |
| `MYSQL_HOST` | store.mysql.host |
| `MYSQL_PORT` | store.mysql.port |
| `MYSQL_NAME` | store.mysql.name |
| `REDIS_HOST` | store.redis.host |
| `REDIS_AUTH` | store.redis.auth |
| `REDIS_NUM` | store.redis.num |
| `STARGO_LOG_LEVEL` | log.level |

## Loading config

```go
import "github.com/starfork/stargo/config"

// From YAML file (default: -c ../config/debug.yaml)
conf, err := config.LoadConfig()

// With explicit path
conf, err := config.LoadConfig("path/to/config.yaml")

// Parse from an existing YAML file
conf, err := config.ParseConfig("path/to/config.yaml")
```
