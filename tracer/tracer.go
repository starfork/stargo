package tracer

import "go.opentelemetry.io/otel/trace"

type Tracer interface {
	Tracer() trace.Tracer
	Close() error
}

var DefaultTracer Tracer = &NoopTracer{}

type NoopTracer struct{}

func (t *NoopTracer) Tracer() trace.Tracer { return trace.NewNoopTracerProvider().Tracer("") }
func (t *NoopTracer) Close() error         { return nil }

var tracerFactories = make(map[string]func(*Config) (Tracer, error))

func Register(name string, factory func(*Config) (Tracer, error)) {
	tracerFactories[name] = factory
}

func NewTracer(name string, conf *Config) (Tracer, error) {
	if f, ok := tracerFactories[name]; ok {
		return f(conf)
	}
	return &NoopTracer{}, nil
}
