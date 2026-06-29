package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	// By default, tracer is a noop — all trace calls are safe.
	// 默认 tracer 为 noop，所有追踪调用均安全。
	//
	// To enable Jaeger tracing:
	// 1. import _ "github.com/starfork/stargo/tracer/jaeger"
	// 2. Configure config.yaml:
	//    tracer:
	//      driver: jaeger
	//      host: "127.0.0.1:6831"
	//      name: "trace-demo"
	// 要启用 Jaeger 追踪:
	// 1. import _ "github.com/starfork/stargo/tracer/jaeger"
	// 2. 在 config.yaml 中配置:
	//    tracer:
	//      driver: jaeger
	//      host: "127.0.0.1:6831"
	//      name: "trace-demo"
	conf, _ := config.LoadConfig()
	app := stargo.New("trace-demo", conf)
	h := NewHandler(app.Logger())

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
