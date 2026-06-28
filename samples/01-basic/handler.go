package main

import (
	"context"

	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 处理器结构体: 嵌入 UnimplementedSampleServiceServer 确保向前兼容
// Handler struct: embeds UnimplementedSampleServiceServer for forward compatibility
type handler struct {
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

// NewHandler 构造函数 — 注入依赖 (logger)
// NewHandler constructor — injects dependencies (logger)
func NewHandler(l logger.Logger) *handler {
	return &handler{
		logger: l,
	}
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser called: id=%d", req.Id)
	// In a real service, query the database here
	return &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser called: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{
		Users: []*pb.GetUserResponse{
			{Id: 1, Name: "Alice", Email: "alice@example.com"},
			{Id: 2, Name: "Bob", Email: "bob@example.com"},
		},
		Total: 2,
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser called: name=%s email=%s", req.Name, req.Email)
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	// In a real service, insert into the database here
	return &pb.CreateUserResponse{Id: 100}, nil
}
