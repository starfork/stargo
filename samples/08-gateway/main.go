package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/api"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("gateway-demo", conf)
	h := NewHandler(app.Logger())

	// 注册 gRPC 服务端 / Register the gRPC server
	pb.RegisterSampleServiceServer(app.RpcServer(), h)

	// 在 goroutine 中启动 gRPC 服务 / Start gRPC in background
	go app.Run(&pb.SampleService_ServiceDesc, h)

	// 启动 HTTP 网关 / Start HTTP gateway
	gw := api.NewApi(&api.Config{
		App:  "gateway-demo",
		Port: ":8080",
		// 可选：加密 marshaler (AES-GCM) / Optional encrypted marshaler
		// Enc:    true,
		// EncKey: "0123456789abcdef0123456789abcdef",
	})
	gw.Run()
}
