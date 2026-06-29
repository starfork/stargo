package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("naming-demo", conf)
	h := NewHandler(app)

	// The service is auto-registered with etcd during Run(),
	// and auto-deregistered during Stop().
	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
