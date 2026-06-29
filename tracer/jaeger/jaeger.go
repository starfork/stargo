package jaeger

import (
	"fmt"
	"io"

	"github.com/starfork/stargo/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func init() {
	tracer.Register("jaeger", NewJaeger)
}

type JaegerTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func (j *JaegerTracer) Close() error {
	return j.closer.Close()
}

func NewJaeger(conf *tracer.Config) (tracer.Tracer, error) {
	cfg := &config.Configuration{
		ServiceName: conf.Name,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	if conf.Host != "" {
		cfg.Reporter.LocalAgentHostPort = conf.Host
	}
	t, c, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, fmt.Errorf("jaeger init: %w", err)
	}
	return &JaegerTracer{tracer: t, closer: c}, nil
}
