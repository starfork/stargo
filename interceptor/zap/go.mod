module github.com/starfork/stargo/interceptor/zap

go 1.26.4

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.4
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.84.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260921155816-b14227669459 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/starfork/stargo => ../../
