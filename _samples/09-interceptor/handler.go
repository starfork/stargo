package main

import (
	"context"

	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(l logger.Logger) *handler {
	return &handler{logger: l}
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)
	return &pb.GetUserResponse{Id: req.Id, Name: "Alice", Email: "alice@example.com"}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}

// SimulatePanic 演示 recovery 拦截器效果 / Demonstrates the recovery interceptor
// 注意: 这个方法不在 proto 定义中, 不会被 gRPC 调用; 仅作为文档示例
// Note: this method is NOT in the proto definition, never called by gRPC; documentation only
func (h *handler) SimulatePanic(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	panic("simulated panic for recovery test")
}
