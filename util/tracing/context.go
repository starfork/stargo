package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func SpanIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}
