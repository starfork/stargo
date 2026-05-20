# 配置参考

[English](../en/config.md)

## YAML 结构

```yaml
env: dev                   # dev | docker | production
timezone: Asia/Shanghai    # 默认时区
timeformat: "2006-01-02T15:04:05+08:00"

server:
  addr: ":9090"
  # UnaryInterceptor:     # 通过代码配置
  # StreamInterceptor:    # 通过代码配置
  # ServerOpts:           # 通过代码配置

api:
  app: my-service
  port: ":8080"
  registry:               # 网关使用的 etcd 配置
    scheme: etcd
    host: localhost:2379
  enc: false              # 是否开启 AES-GCM 加密
  enckey: ""              # 加密密钥（32字节）

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
    num: 0                # 数据库编号
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
  ttl: 10                 # 租约 TTL（秒）

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

## 环境变量覆盖

| 变量 | 覆盖配置 |
|------|----------|
| `MYSQL_USER` | store.mysql.user |
| `MYSQL_PASSWD` | store.mysql.auth |
| `MYSQL_HOST` | store.mysql.host |
| `MYSQL_PORT` | store.mysql.port |
| `MYSQL_NAME` | store.mysql.name |
| `REDIS_HOST` | store.redis.host |
| `REDIS_AUTH` | store.redis.auth |
| `REDIS_NUM` | store.redis.num |
| `STARGO_LOG_LEVEL` | log.level |

## 加载配置

```go
import "github.com/starfork/stargo/config"

// 从 YAML 文件加载（默认: -c ../config/debug.yaml）
conf, err := config.LoadConfig()

// 指定路径
conf, err := config.LoadConfig("path/to/config.yaml")

// 解析已有 YAML 文件
conf, err := config.ParseConfig("path/to/config.yaml")
```
