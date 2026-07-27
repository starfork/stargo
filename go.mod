module github.com/starfork/stargo

go 1.26.4

retract (
	[v0.1.1, v0.1.9]
	[v0.0.1, v0.0.8]
)

// bugs found, not support

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/prometheus/client_golang v1.24.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.69.0
	go.opentelemetry.io/otel/trace v1.44.0
	golang.org/x/exp v0.0.0-20260718201538-764159d718ef
	golang.org/x/sync v0.22.0
	golang.org/x/text v0.40.0
	golang.org/x/time v0.15.0
	google.golang.org/grpc v1.82.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260724162435-b2f20204f0df // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260724162435-b2f20204f0df
	google.golang.org/protobuf v1.36.11
)
