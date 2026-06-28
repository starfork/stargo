package server

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcServerHandlingSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_server_handling_seconds",
		Help:    "Histogram of gRPC server handling time in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"grpc_service", "grpc_method"})

	grpcServerHandledTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_server_handled_total",
		Help: "Total number of gRPC requests completed.",
	}, []string{"grpc_service", "grpc_method", "grpc_code"})

	grpcServerStartedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_server_started_total",
		Help: "Total number of gRPC requests started.",
	}, []string{"grpc_service", "grpc_method"})
)

func UnaryServerMetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	service, method := splitMethod(info.FullMethod)
	grpcServerStartedTotal.WithLabelValues(service, method).Inc()

	start := time.Now()
	resp, err := handler(ctx, req)
	elapsed := time.Since(start)

	code := status.Code(err)
	grpcServerHandledTotal.WithLabelValues(service, method, code.String()).Inc()
	grpcServerHandlingSeconds.WithLabelValues(service, method).Observe(elapsed.Seconds())

	return resp, err
}

func StreamServerMetricsInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	service, method := splitMethod(info.FullMethod)
	grpcServerStartedTotal.WithLabelValues(service, method).Inc()

	start := time.Now()
	err := handler(srv, ss)
	elapsed := time.Since(start)

	code := status.Code(err)
	grpcServerHandledTotal.WithLabelValues(service, method, code.String()).Inc()
	grpcServerHandlingSeconds.WithLabelValues(service, method).Observe(elapsed.Seconds())

	return err
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func splitMethod(fullMethod string) (string, string) {
	if len(fullMethod) == 0 {
		return "", ""
	}
	// format: /package.Service/Method
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
