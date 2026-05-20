package main

import (
	"context"
	"log"

	_ "github.com/starfork/stargo/store/mysql"
	_ "github.com/starfork/stargo/store/redis"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	smysql "github.com/starfork/stargo/store/mysql"
	sredis "github.com/starfork/stargo/store/redis"
	"gorm.io/gorm"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app := stargo.New("store-demo", conf)
	ctx := context.Background()

	// MySQL access (requires blank import of store/mysql and YAML config)
	if db := app.Store("mysql"); db != nil {
		gormDB := db.(*smysql.Mysql).Instance().(*gorm.DB)
		var version string
		if err := gormDB.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
			app.LogErrorf("mysql query: %v", err)
		} else {
			app.LogInfof("mysql version: %s", version)
		}
	} else {
		app.LogInfof("mysql not configured")
	}

	// Redis access (requires blank import of store/redis and YAML config)
	if rdc := app.Store("redis"); rdc != nil {
		client := rdc.(*sredis.Redis).Instance().(*redis.Client)
		if err := client.Set(ctx, "demo:key", "hello stargo", 0).Err(); err != nil {
			app.LogErrorf("redis set: %v", err)
		} else {
			val, _ := client.Get(ctx, "demo:key").Result()
			app.LogInfof("redis value: %s", val)
		}
	} else {
		app.LogInfof("redis not configured")
	}
}
