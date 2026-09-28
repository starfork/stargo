module github.com/starfork/stargo/store/redis

go 1.26.4

require (
	github.com/redis/go-redis/v9 v9.22.0
	github.com/starfork/stargo v1.1.10
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

replace github.com/starfork/stargo => ../../
