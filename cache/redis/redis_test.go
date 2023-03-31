package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/config"
	sredis "github.com/starfork/stargo/store/redis"
)

var rdc *redis.Client

var test_conf = &config.ServerConfig{
	Redis: &config.RedisConfig{
		Addr: "127.0.0.1:6379",
	},
}

func TestGet(t *testing.T) {

	r := sredis.Connect(test_conf)
	rdc = r.GetInstance()

	type Stat struct {
		Total  uint32  `json:"total"`
		Amount float64 `json:"amount"`
	}
	stat := &Stat{
		Total:  1000,
		Amount: 93.25,
	}
	key := "sdfsfsdfsdfsd"
	ctx := context.Background()
	rdc.SetNX(ctx, key, stat, time.Second*6000)
	data := &Stat{}
	rdc.Get(ctx, key).Scan(&data)
	fmt.Printf("%+v", data)

}
