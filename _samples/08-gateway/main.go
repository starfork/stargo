package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/api"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("gateway-demo", conf)
	h := NewHandler(app.Logger())

	pb.RegisterSampleServiceServer(app.RpcServer(), h)

	go app.Run(&pb.SampleService_ServiceDesc, h)

	gw, err := api.NewApi(&api.Config{
		App:      "gateway-demo",
		Port:     ":8080",
		Registry: conf.Registry,
		DiaOpts:  []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	})
	api.E(err)

	api.E(pb.RegisterSampleServiceHandlerClient(gw.Ctx(), gw.Rmux(), pb.NewSampleServiceClient(gw.Conn())))

	api.E(gw.Run())
}
