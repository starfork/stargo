module github.com/starfork/stargo

go 1.26.2

retract (
	[v0.1.1, v0.1.9]
	[v0.0.1, v0.0.8]
)

// bugs found, not support

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b
	golang.org/x/text v0.30.0
	golang.org/x/time v0.14.0
	google.golang.org/grpc v1.76.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff
	google.golang.org/protobuf v1.36.10
)
