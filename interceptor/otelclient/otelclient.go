package otelclient

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func New() stats.Handler {
	return otelgrpc.NewClientHandler()
}

func DialOption() grpc.DialOption {
	return grpc.WithStatsHandler(New())
}
