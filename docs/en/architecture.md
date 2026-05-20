# Architecture Overview

[中文](../zh/architecture.md)

## Module structure

```
stargo/
├── app.go              # App struct, lifecycle (New → Run → Stop)
├── stargo.go           # Accessor methods (RpcServer, Config, Store, ...)
├── options.go          # Functional options pattern
├── logger.go           # Logger convenience methods on App
│
├── server/             # gRPC server wrapper
│   └── server.go       #   New → Run → Stop/Restart
│
├── config/             # YAML config loading + etcd config manager
│
├── client/             # gRPC client with service discovery
│
├── api/                # HTTP gateway (grpc-gateway)
│   ├── api.go          #   Gateway setup + Run
│   ├── cors.go         #   CORS middleware
│   ├── meta.go         #   Metadata extraction helpers
│   └── custom/         #   AES-GCM encrypted marshaler
│
├── broker/             # Message broker interface
│   └── nats/           #   NATS implementation
│
├── naming/             # Service registry & resolver interfaces
│   └── etcd/           #   etcd implementation
│
├── store/              # Storage interface
│   ├── mysql/          #   MySQL/GORM implementation
│   └── redis/          #   Redis implementation
│
├── tracer/             # Tracing interface
│   └── jaeger/         #   Jaeger/OpenTracing implementation
│
├── queue/              # Delayed task queue (Redis sorted sets)
│
├── interceptor/        # gRPC interceptors
│   ├── auth/           #   Authentication
│   ├── logger/zap/     #   Zap logging
│   ├── ratelimit/      #   Rate limiting
│   ├── recovery/       #   Panic recovery
│   └── validator/      #   Struct validation
│
├── cache/              # Cache abstraction
│   └── filecache/      #   File-based implementation
│
├── filemanager/        # File storage interface
│
├── pm/                 # Parameter map (typed getters, URL encoding)
│
└── util/               # Utility packages
    ├── ustring/        #   String manipulation
    ├── request/        #   HTTP helpers
    ├── geo/            #   Geography (Haversine)
    ├── merger/         #   Data merging
    ├── number/         #   Number utilities
    ├── parser/         #   ID card parser
    ├── path/           #   User path encoding
    └── urange/         #   Reflection-based field iteration
```

## Lifecycle

```
stargo.New("name", conf)
  └─ initConfig()
       ├─ initialize stores (if configured + explicitly imported)
       ├─ initialize NATS broker (if configured)
       ├─ initialize etcd registry (if configured)
       └─ set default tracer

app.Run(desc, impl)
  └─ beforeRun()
  │    ├─ set timezone
  │    ├─ create gRPC server (with interceptors)
  │    ├─ register reflection (non-production only)
  │    └─ register service with etcd (if configured)
  ├─ register services on gRPC server
  ├─ signal handler (SIGTERM/SIGINT/SIGHUP/SIGQUIT)
  └─ server.Run()
       └─ gRPC Serve()

app.Stop()
  ├─ deregister from etcd
  ├─ close all stores
  ├─ unsubscribe broker
  └─ close tracer
```

## Store opt-in

Stores use a registry pattern: each store package registers itself via `init()`. The root `app.go` does NOT import mysql/redis packages. Users must blank-import the stores they need:

```go
import _ "github.com/starfork/stargo/store/mysql"
import _ "github.com/starfork/stargo/store/redis"
```

This keeps the dependency tree minimal by default.
