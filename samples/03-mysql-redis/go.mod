module github.com/starfork/stargo/samples/03-mysql-redis

go 1.26.2

require (
	github.com/redis/go-redis/v9 v9.14.0
	github.com/starfork/stargo v0.0.0
	github.com/starfork/stargo/samples/proto/sample v0.0.0
	google.golang.org/grpc v1.76.0
	gorm.io/gorm v1.31.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/starfork/stargo => ../../
	github.com/starfork/stargo/samples/proto/sample => ../proto/sample
)
