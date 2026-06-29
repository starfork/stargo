package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("logger-demo", conf)
	h := NewHandler(app.Logger())

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
