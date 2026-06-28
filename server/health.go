package server

import (
	"context"
	"sync"

	"github.com/starfork/stargo/logger"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	mu     sync.RWMutex
	status map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		status: map[string]grpc_health_v1.HealthCheckResponse_ServingStatus{
			"": grpc_health_v1.HealthCheckResponse_SERVING,
		},
	}
}

func (h *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	service := req.GetService()
	status, ok := h.status[service]
	if !ok {
		logger.DefaultLogger.Warnf("health check for unknown service: %s", service)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &grpc_health_v1.HealthCheckResponse{Status: status}, nil
}

func (h *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	service := req.GetService()
	h.mu.RLock()
	status, ok := h.status[service]
	h.mu.RUnlock()
	if !ok {
		status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	if err := stream.Send(&grpc_health_v1.HealthCheckResponse{Status: status}); err != nil {
		return err
	}
	return nil
}

func (h *HealthServer) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status[service] = status
}
