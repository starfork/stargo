package tracer

type Config struct {
	Driver string // tracer driver: "jaeger", "otel"
	Host   string
	Name   string // service name
}
