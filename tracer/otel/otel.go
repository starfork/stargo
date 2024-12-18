package otel

import (
	"context"

	"github.com/starfork/stargo/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Otel struct {
	ctx context.Context
	p   *sdktrace.TracerProvider
}

func NewTracer(conf *tracer.Config) (tracer.Tracer, error) {
	ctx := context.Background()
	o := &Otel{ctx: ctx}
	// 配置 OTLP gRPC Exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(conf.Host), // Jaeger 的 OTLP 接收器地址
		otlptracegrpc.WithInsecure(),          // 禁用 TLS（开发环境使用）
	)
	if err != nil {
		return nil, err
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // 使用批量发送器
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.Name), // 服务名称
		)),
	)
	otel.SetTracerProvider(tp)
	o.p = tp
	return o, nil
}

func (p *Otel) Close() error {
	return p.p.Shutdown(context.Background())
}
