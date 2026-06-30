module github.com/starfork/stargo/cache/redis

go 1.26.4

require (
	github.com/json-iterator/go v1.1.12
	github.com/redis/go-redis/v9 v9.21.0
	github.com/starfork/stargo v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace github.com/starfork/stargo => ../../
