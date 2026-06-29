module github.com/starfork/stargo/store/mysql

go 1.26.2

require (
	github.com/redis/go-redis/v9 v9.19.0
	github.com/starfork/gorm-cache v0.0.0-20251013074659-4bf32fdac72c
	github.com/starfork/stargo v0.0.0
	github.com/twpayne/go-geom v1.6.1
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

replace github.com/starfork/stargo => ../../
