package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	_ "github.com/starfork/stargo/store/mysql" // MySQL 存储驱动 / MySQL store driver
	_ "github.com/starfork/stargo/store/redis" // Redis 存储驱动 / Redis store driver
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("store-demo", conf)
	h := NewHandler(app)

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
