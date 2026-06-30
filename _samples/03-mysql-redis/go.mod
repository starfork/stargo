module github.com/starfork/stargo/samples/03-mysql-redis

go 1.26.4

require (
	github.com/redis/go-redis/v9 v9.21.0
	github.com/starfork/stargo v0.0.0
	github.com/starfork/stargo/samples/proto/sample v0.0.0
	github.com/starfork/stargo/store/mysql v0.0.0
	github.com/starfork/stargo/store/redis v0.0.0
	google.golang.org/grpc v1.76.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/starfork/gorm-cache v0.0.0-20251013074659-4bf32fdac72c // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
)

replace (
	github.com/starfork/stargo => ../../
	github.com/starfork/stargo/samples/proto/sample => ../proto/sample
	github.com/starfork/stargo/store/mysql => ../../store/mysql
	github.com/starfork/stargo/store/redis => ../../store/redis
)
