package otelserver

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func New() stats.Handler {
	return otelgrpc.NewServerHandler()
}

func StatsHandlerOption() grpc.ServerOption {
	return grpc.StatsHandler(New())
}
