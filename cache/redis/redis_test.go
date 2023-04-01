package redis

import (
	"context"
	"encoding"
	"fmt"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
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

type TestStat struct {
	Total  uint32  `json:"total"`
	Amount float64 `json:"amount"`
}

var _ encoding.BinaryMarshaler = new(TestStat)
var _ encoding.BinaryUnmarshaler = new(TestStat)

func (e *TestStat) MarshalBinary() (data []byte, err error) {
	return jsoniter.Marshal(e)
}
func (e *TestStat) UnmarshalBinary(data []byte) error {

	return jsoniter.Unmarshal(data, e)
}

func TestGet(t *testing.T) {

	r := sredis.Connect(test_conf)
	rdc = r.GetInstance()

	stat := &TestStat{
		Total:  1000,
		Amount: 93.25,
	}
	key := "sdfsfsdfsdfsd"
	ctx := context.Background()
	rdc.SetNX(ctx, key, stat, time.Second*6000)
	data := &TestStat{}

	if err := rdc.Get(ctx, key).Scan(data); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", data)

}
