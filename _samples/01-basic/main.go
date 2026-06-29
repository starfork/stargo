package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	// Config-first: stargo reads YAML and auto-connects configured components.
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := stargo.New("basic-service", conf)
	h := NewHandler(app.Logger())

	// Register the gRPC service and start serving.
	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
