package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	// By default, tracer.DefaultTracer is a noop — all trace calls are safe.
	//
	// To enable Jaeger tracing, import the Jaeger tracer and set it before New:
	//
	//   import jtracer "github.com/starfork/stargo/tracer/jaeger"
	//   tracer.DefaultTracer = jtracer.InitJaeger("trace-demo")
	//
	conf, _ := config.LoadConfig()
	app := stargo.New("trace-demo", conf)
	h := NewHandler(app.Logger())

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
