module github.com/starfork/stargo/samples/04-trace

go 1.26.2

require (
	github.com/starfork/stargo v0.0.0
	github.com/starfork/stargo/samples/proto/sample v0.0.0
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/grpc v1.76.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/starfork/stargo => ../../
	github.com/starfork/stargo/samples/proto/sample => ../proto/sample
)
