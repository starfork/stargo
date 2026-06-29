# stargo 全栈微服务示例 / Full-stack Microservices Demo

4 个微服务: user-service → product-service → order-service → gateway，覆盖 stargo 框架 90% 的功能点。

## 架构 / Architecture

```
                    ┌─────────────┐
                    │   Gateway   │  HTTP :8080 (对外 / external)
                    │  (grpc-gw)  │  gRPC :9099 (内部 / internal)
                    └──────┬──────┘
                           │ etcd 服务发现 / etcd discovery
              ┌────────────┼────────────┬────────────────┐
              ▼            ▼            ▼                ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐
    │ user-service│ │product-svc  │ │order-service│ │  Jaeger   │
    │   :9091     │ │   :9093     │ │   :9092     │ │  :16686   │
    └──┬──────┬───┘ └──────┬──────┘ └──┬──────────┘ └───────────┘
       │      │            │           │
       ▼      ▼            ▼           ▼
    ┌──────┐ ┌──────┐  ┌──────┐   ┌──────┐
    │MySQL │ │Redis │  │MySQL │   │ etcd │
    │:3306 │ │:6379 │  │:3306 │   │:2379 │
    └──────┘ └──────┘  └──────┘   └──────┘
```

### 调用链 / Call Chain
```
Client → Gateway(:8080)
           │
           ├── GET  /api/users/{id}       → etcd → user-service(:9091)
           │                                                ├── Redis cache
           │                                                └── MySQL users
           │
           ├── GET  /api/products/{id}    → etcd → product-service(:9093)
           │                                                └── MySQL products
           │
           └── POST /api/orders           → etcd → order-service(:9092)
                                                      ├── call user-service.GetUser()       ← 验证用户
                                                      ├── call product-service.GetProduct() ← 验证商品+价格
                                                      └── MySQL orders
```
链路追踪 / Tracing:
```
gateway → order-service → user-service (同一 trace-id)
                        → product-service (同一 trace-id)
```

## 功能清单 / Feature Checklist

| 功能 / Feature | user-service | product-service | order-service | gateway |
|---|---|---|---|---|
| MySQL (GORM) | ✅ 用户表 CRUD | ✅ 商品表 CRUD | ✅ 订单表 CRUD | - |
| Redis 缓存 | ✅ Cache-Aside | - | - | - |
| Auth 拦截器 | ✅ Bearer token | - | - | - |
| Rate Limit 拦截器 | - | - | ✅ 10 req/s | - |
| etcd 服务注册 | ✅ | ✅ | ✅ | ✅ |
| etcd 服务发现 (resolver) | - | - | ✅ 发现 user+product | ✅ 发现后端 |
| gRPC Client 调用 | - | - | ✅ → user & product | ✅ → 后端 |
| Zap 结构化日志 | ✅ JSON | ✅ JSON | ✅ JSON | ✅ JSON |
| Jaeger 链路追踪 | ✅ | ✅ | ✅ | - |
| Prometheus metrics | ✅ /metrics | ✅ /metrics | ✅ /metrics | ✅ /metrics |

## 启动 / Quick Start

### 方式 1: Docker Compose (推荐 / Recommended)
```bash
cd _samples/14-full-stack
docker compose up -d

# 查看日志
docker compose logs -f order-service

# 查看 Jaeger 界面
open http://localhost:16686

# 测试
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics

# 停止
docker compose down
```

### 方式 2: 本地开发 / Local Development
```bash
# 1. 启动基础设施
docker compose up -d mysql redis etcd jaeger

# 2. 等待 MySQL 就绪
until docker exec stargo-mysql mysqladmin ping -h localhost --silent; do sleep 1; done

# 3. 初始化表结构
docker exec -i stargo-mysql mysql -uroot -proot123 stargo_demo < init.sql

# 4. 分别启动服务 (四个终端 / four terminals)
cd user-service    && go run .    # 终端 1: user-service     :9091
cd product-service && go run .    # 终端 2: product-service  :9093
cd order-service   && go run .    # 终端 3: order-service    :9092
cd gateway         && go run .    # 终端 4: gateway          :8080
```

### 方式 3: K8s 部署 / K8s Deploy
```bash
kubectl apply -f k8s/
kubectl -n stargo-demo get pods -w
```

## 测试 / Testing

```bash
# 健康检查 / Health check
curl http://localhost:8080/healthz

# 创建用户 / Create user (requires auth)
curl -X POST http://localhost:9091/user.v1.UserService/CreateUser \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d '{"name":"David","email":"david@test.com"}'

# 查询用户 / Query user (verifies Redis cache)
curl -X POST http://localhost:9091/user.v1.UserService/GetUser \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d '{"id":1}'

# 查询商品 / Query product
curl -X POST http://localhost:9093/product.v1.ProductService/GetProduct \
  -H "Content-Type: application/json" \
  -d '{"id":1}'

# 创建订单 / Create order (cross-service: validates user + product)
curl -X POST http://localhost:9092/order.v1.OrderService/CreateOrder \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product":"iPad","amount":3499}'

# Prometheus 指标 / Prometheus metrics
curl http://localhost:8080/metrics
```

## 配置说明 / Config Notes

| 环境变量 / Env Var | 用途 / Purpose |
|---|---|
| `MYSQL_HOST` | MySQL 地址 (替代 config.yaml) |
| `MYSQL_USER` / `MYSQL_PASSWD` | MySQL 认证 |
| `MYSQL_NAME` | 数据库名 |
| `REDIS_HOST` | Redis 地址 |
| `STARGO_LOG_LEVEL` | 日志级别 (trace/debug/info/warn/error) |
