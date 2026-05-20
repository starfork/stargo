package main

import (
	"context"
	"time"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	app    *stargo.App
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	return &handler{
		app:    app,
		logger: app.Logger(),
	}
}

// GetUser calls a downstream service via the gRPC client.
func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	// app.Client() requires a configured registry/resolver.
	cli := h.app.Client()
	if cli == nil {
		return nil, status.Error(codes.Unavailable, "service discovery not configured")
	}

	// Connect to the downstream "user-service".
	conn, err := cli.NewClient("user-service",
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "connect to user-service: %v", err)
	}
	defer conn.Close()

	// Use conn to create a gRPC client for the downstream service:
	// downstreamClient := pb.NewUserServiceClient(conn)
	_ = conn

	return &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	// Call another service's gRPC endpoint.
	cli := h.app.Client()
	if cli != nil {
		conn, err := cli.NewClient("user-service")
		if err == nil {
			defer conn.Close()
			// downstreamClient := pb.NewUserServiceClient(conn)
			// resp, err := downstreamClient.CreateUser(ctx, req)
		}
	}

	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
