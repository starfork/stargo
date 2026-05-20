package main

import (
	"context"

	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
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

// Demonstrate different log levels inside a service method.
func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Debugf("GetUser request received: id=%d", req.Id)

	if req.Id <= 0 {
		h.logger.Warnf("GetUser called with invalid id: %d", req.Id)
		return nil, nil
	}

	// In a real service, query logic goes here.
	h.logger.Infof("GetUser success: id=%d", req.Id)

	return &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	if req.Email == "" {
		h.logger.Errorf("CreateUser missing email for user: %s", req.Name)
		return nil, nil
	}

	h.logger.Infof("CreateUser success: name=%s", req.Name)
	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
