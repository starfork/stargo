module github.com/starfork/stargo/interceptor/zap

go 1.26.4

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.4
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.83.2
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260911204522-f61a6ca850bd // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/starfork/stargo => ../../
