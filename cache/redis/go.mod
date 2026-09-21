module github.com/starfork/stargo/cache/redis

go 1.26.4

require (
	github.com/json-iterator/go v1.1.12
	github.com/redis/go-redis/v9 v9.22.0
	github.com/starfork/stargo v1.1.9
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
)

replace github.com/starfork/stargo => ../../
