package main

import (
	"context"

	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"github.com/starfork/stargo/tracer"
)

type handler struct {
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(l logger.Logger) *handler {
	return &handler{
		logger: l,
	}
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	// stargo.DefaultTracer is a noop. Replace it with a real implementation
	// (e.g., Jaeger) before stargo.New to enable distributed tracing.
	_ = tracer.DefaultTracer

	return &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)
	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
