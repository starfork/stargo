package client

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcClientHandlingSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_client_handling_seconds",
		Help:    "Histogram of gRPC client handling time in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"grpc_service", "grpc_method"})

	grpcClientHandledTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_client_handled_total",
		Help: "Total number of gRPC client requests completed.",
	}, []string{"grpc_service", "grpc_method", "grpc_code"})

	grpcClientStartedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_client_started_total",
		Help: "Total number of gRPC client requests started.",
	}, []string{"grpc_service", "grpc_method"})

	clientRegistry = prometheus.NewRegistry()
)

func init() {
	clientRegistry.MustRegister(grpcClientHandlingSeconds)
	clientRegistry.MustRegister(grpcClientHandledTotal)
	clientRegistry.MustRegister(grpcClientStartedTotal)
}

func ClientRegistry() *prometheus.Registry {
	return clientRegistry
}

func UnaryClientMetricsInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	service, mtd := splitClientMethod(method)
	grpcClientStartedTotal.WithLabelValues(service, mtd).Inc()

	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)

	code := status.Code(err)
	grpcClientHandledTotal.WithLabelValues(service, mtd, code.String()).Inc()
	grpcClientHandlingSeconds.WithLabelValues(service, mtd).Observe(elapsed.Seconds())

	return err
}

func StreamClientMetricsInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	service, mtd := splitClientMethod(method)
	grpcClientStartedTotal.WithLabelValues(service, mtd).Inc()

	start := time.Now()
	cs, err := streamer(ctx, desc, cc, method, opts...)
	elapsed := time.Since(start)

	code := status.Code(err)
	grpcClientHandledTotal.WithLabelValues(service, mtd, code.String()).Inc()
	grpcClientHandlingSeconds.WithLabelValues(service, mtd).Observe(elapsed.Seconds())

	return cs, err
}

func splitClientMethod(fullMethod string) (string, string) {
	if len(fullMethod) == 0 {
		return "", ""
	}
	service := ""
	method := ""
	for i := 1; i < len(fullMethod); i++ {
		if fullMethod[i] == '/' {
			service = fullMethod[1:i]
			method = fullMethod[i+1:]
			break
		}
	}
	return service, method
}
