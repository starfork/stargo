package tracer

type Tracer interface {
	Close() error
}

var DefaultTracer Tracer = &NoopTracer{}

type NoopTracer struct{}

func (t *NoopTracer) Close() error { return nil }
