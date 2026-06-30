package server

import (
	"context"
	"sync"

	"github.com/starfork/stargo/logger"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type subscriber struct {
	ch     chan grpc_health_v1.HealthCheckResponse_ServingStatus
	cancel context.CancelFunc
}

type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	mu          sync.RWMutex
	status      map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
	subscribers map[string][]*subscriber

	// Readiness tracking
	deps   map[string]bool // dependency name -> healthy
	depsMu sync.RWMutex
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		status: map[string]grpc_health_v1.HealthCheckResponse_ServingStatus{
			"": grpc_health_v1.HealthCheckResponse_SERVING,
		},
		subscribers: make(map[string][]*subscriber),
		deps:        make(map[string]bool),
	}
}

func (h *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	service := req.GetService()
	// Readiness check: aggregate dependency health
	if service == "readiness" {
		status := h.aggregateReadiness()
		return &grpc_health_v1.HealthCheckResponse{
			Status: status,
		}, nil
	}

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
	
	// Create context for this watch stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	
	// Create subscriber
	sub := &subscriber{
		ch:     make(chan grpc_health_v1.HealthCheckResponse_ServingStatus, 1),
		cancel: cancel,
	}
	
	// Register subscriber
	h.mu.Lock()
	h.subscribers[service] = append(h.subscribers[service], sub)
	// Get current status
	status, ok := h.status[service]
	h.mu.Unlock()
	
	if !ok {
		status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	
	// Send initial status
	if err := stream.Send(&grpc_health_v1.HealthCheckResponse{Status: status}); err != nil {
		h.removeSubscriber(service, sub)
		return err
	}
	
	// Wait for status updates or context cancellation
	for {
		select {
		case <-ctx.Done():
			h.removeSubscriber(service, sub)
			return nil
		case newStatus := <-sub.ch:
			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{Status: newStatus}); err != nil {
				h.removeSubscriber(service, sub)
				return err
			}
		}
	}
}

func (h *HealthServer) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status[service] = status
	
	// Notify subscribers
	for _, sub := range h.subscribers[service] {
		select {
		case sub.ch <- status:
		default:
			// Channel full, skip this subscriber
		}
	}
}

func (h *HealthServer) removeSubscriber(service string, sub *subscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()

	subs := h.subscribers[service]
	for i, s := range subs {
		if s == sub {
			h.subscribers[service] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// SetDependency marks a named dependency as healthy or unhealthy
func (h *HealthServer) SetDependency(name string, healthy bool) {
	h.depsMu.Lock()
	h.deps[name] = healthy
	h.depsMu.Unlock()
}

// aggregateReadiness returns SERVING only if all dependencies are healthy
func (h *HealthServer) aggregateReadiness() grpc_health_v1.HealthCheckResponse_ServingStatus {
	h.depsMu.RLock()
	defer h.depsMu.RUnlock()

	// Require at least one dependency tracked
	if len(h.deps) == 0 {
		// No dependencies configured, return SERVING (best effort)
		return grpc_health_v1.HealthCheckResponse_SERVING
	}

	for _, healthy := range h.deps {
		if !healthy {
			return grpc_health_v1.HealthCheckResponse_NOT_SERVING
		}
	}
	return grpc_health_v1.HealthCheckResponse_SERVING
}
