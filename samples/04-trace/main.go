package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/tracer"
)

func main() {
	// stargo uses a noop tracer (tracer.DefaultTracer) by default.
	// All trace operations are safe to call even with the noop tracer.
	//
	// To enable Jaeger tracing, import and use the Jaeger tracer:
	//
	//   import jtracer "github.com/starfork/stargo/tracer/jaeger"
	//   otTracer, closer := jtracer.InitJaeger("trace-demo")
	//   _ = otTracer
	//   _ = closer

	conf, _ := config.LoadConfig()
	app := stargo.New("trace-demo", conf)

	// Verify the tracer is a noop by default
	_ = tracer.DefaultTracer
	app.LogInfof("tracer type: %T", app.Tracer())

	// Use the tracer interface (Close is called automatically on app.Stop())
	if err := app.Tracer().Close(); err != nil {
		app.LogErrorf("tracer close: %v", err)
	}
}
