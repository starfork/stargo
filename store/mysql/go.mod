module github.com/starfork/stargo/store/mysql

go 1.26.4

require (
	github.com/redis/go-redis/v9 v9.23.0
	github.com/starfork/gorm-cache v0.0.0-20251013074659-4bf32fdac72c
	github.com/starfork/stargo v1.1.12
	github.com/twpayne/go-geom v1.7.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.2
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-sql-driver/mysql v1.10.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace github.com/starfork/stargo => ../../
