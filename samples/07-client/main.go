package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("client-demo", conf)
	h := NewHandler(app)

	// The gRPC client uses the configured resolver for service discovery.
	// When the handler's GetUser or CreateUser is called, it connects
	// to the downstream "user-service" via etcd.
	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
