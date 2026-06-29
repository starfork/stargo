module github.com/starfork/stargo/tracer/jaeger

go 1.26.4

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/starfork/stargo v0.0.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace github.com/starfork/stargo => ../../
