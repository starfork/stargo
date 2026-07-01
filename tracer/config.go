package tracer

import "time"

type Config struct {
	Driver       string        // tracer driver: "jaeger", "otel"
	Host         string        // collector endpoint, e.g. "localhost:4317"
	Name         string        // service name
	SampleRate   float64       // 采样率, 默认 1.0
	BatchTimeout time.Duration // batch timeout, 默认 5s
}
