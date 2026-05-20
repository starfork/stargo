# 架构概览

[English](../en/architecture.md)

## 模块结构

```
stargo/
├── app.go              # App 结构体，生命周期 (New → Run → Stop)
├── stargo.go           # 访问方法 (RpcServer, Config, Store, ...)
├── options.go          # 函数式选项模式
├── logger.go           # App 的日志便捷方法
│
├── server/             # gRPC 服务包装
│   └── server.go       #   New → Run → Stop/Restart
│
├── config/             # YAML 配置加载 + etcd 配置管理
│
├── client/             # 带服务发现的 gRPC 客户端
│
├── api/                # HTTP 网关 (grpc-gateway)
│   ├── api.go          #   网关设置 + Run
│   ├── cors.go         #   CORS 中间件
│   ├── meta.go         #   元数据提取辅助
│   └── custom/         #   AES-GCM 加密 marshaler
│
├── broker/             # 消息代理接口
│   └── nats/           #   NATS 实现
│
├── naming/             # 服务注册与发现接口
│   └── etcd/           #   etcd 实现
│
├── store/              # 存储接口
│   ├── mysql/          #   MySQL/GORM 实现
│   └── redis/          #   Redis 实现
│
├── tracer/             # 链路追踪接口
│   └── jaeger/         #   Jaeger/OpenTracing 实现
│
├── queue/              # 延迟任务队列 (Redis 有序集合)
│
├── interceptor/        # gRPC 拦截器
│   ├── auth/           #   认证
│   ├── logger/zap/     #   Zap 日志
│   ├── ratelimit/      #   限流
│   ├── recovery/       #   panic 恢复
│   └── validator/      #   结构体验证
│
├── cache/              # 缓存抽象
│   └── filecache/      #   文件缓存实现
│
├── filemanager/        # 文件存储接口
│
├── pm/                 # 参数映射（类型安全 getter, URL 编码）
│
└── util/               # 工具包
    ├── ustring/        #   字符串处理
    ├── request/        #   HTTP 辅助
    ├── geo/            #   地理计算 (Haversine)
    ├── merger/         #   数据合并
    ├── number/         #   数字工具
    ├── parser/         #   身份证解析
    ├── path/           #   用户路径编码
    └── urange/         #   基于反射的字段遍历
```

## 生命周期

```
stargo.New("name", conf)
  └─ initConfig()
       ├─ 初始化 stores（如果配置且导入了包）
       ├─ 初始化 NATS broker（如果配置）
       ├─ 初始化 etcd registry（如果配置）
       └─ 设置默认 tracer

app.Run(desc, impl)
  └─ beforeRun()
  │    ├─ 设置时区
  │    ├─ 创建 gRPC 服务（含拦截器）
  │    ├─ 注册 reflection（非生产环境）
  │    └─ 向 etcd 注册服务（如果配置）
  ├─ 在 gRPC 服务上注册服务
  ├─ 信号处理 (SIGTERM/SIGINT/SIGHUP/SIGQUIT)
  └─ server.Run()
       └─ gRPC Serve()

app.Stop()
  ├─ 从 etcd 反注册
  ├─ 关闭所有 stores
  ├─ 取消订阅 broker
  └─ 关闭 tracer
```

## Store opt-in

Stores 使用注册模式：每个 store 包通过 `init()` 自注册。根 `app.go` 不导入 mysql/redis 包。用户必须 blank-import 需要的存储：

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

这样可以保持最小依赖树。
