module github.com/starfork/stargo

go 1.25.2

retract (
	[v0.1.1, v0.1.9]
	[v0.0.1, v0.0.8]
)

// bugs found, not support

require (
	github.com/go-playground/locales v0.14.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2
	github.com/json-iterator/go v1.1.12
	github.com/nats-io/nats.go v1.46.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/redis/go-redis/v9 v9.14.0
	github.com/twpayne/go-geom v1.6.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	go.etcd.io/etcd/client/v3 v3.6.5
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	go.uber.org/ratelimit v0.3.1
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b
	golang.org/x/text v0.30.0
	golang.org/x/time v0.14.0
	google.golang.org/grpc v1.76.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.8.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.28.0
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/starfork/gorm-cache v0.0.0-20251013074659-4bf32fdac72c
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.6.5 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.5 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff
	google.golang.org/protobuf v1.36.10
)
