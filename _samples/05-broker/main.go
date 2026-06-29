package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("broker-demo", conf)
	h := NewHandler(app)

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
