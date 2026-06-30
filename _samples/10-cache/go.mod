module github.com/starfork/stargo/samples/10-cache

go 1.26.4

require (
	github.com/starfork/stargo v0.0.0
	github.com/starfork/stargo/cache/redis v0.0.0
	github.com/starfork/stargo/samples/proto/sample v0.0.0
	github.com/starfork/stargo/store/redis v0.0.0
	google.golang.org/grpc v1.79.3
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/redis/go-redis/v9 v9.21.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/starfork/stargo => ../../
	github.com/starfork/stargo/cache/redis => ../../cache/redis
	github.com/starfork/stargo/samples/proto/sample => ../proto/sample
	github.com/starfork/stargo/store/redis => ../../store/redis
)
