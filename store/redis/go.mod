module github.com/starfork/stargo/store/redis

go 1.26.4

require (
	github.com/redis/go-redis/v9 v9.21.0
	github.com/starfork/stargo v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

replace github.com/starfork/stargo => ../../
