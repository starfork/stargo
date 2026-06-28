package main

import (
	"context"

	"github.com/starfork/stargo/api"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/metadata"
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

	// 从 gRPC 元数据中提取 HTTP 头部 (网关自动将 X-* 等头部转发为 Grpc-Metadata-*)
	// Extract HTTP headers from gRPC metadata (gateway forwards X-* headers as Grpc-Metadata-*)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if token := md.Get("grpc-metadata-x-token"); len(token) > 0 {
			h.logger.Debugf("token from header: %s", token[0])
		}
		if lang := api.MetaLang(ctx); lang != "" {
			h.logger.Debugf("language: %s", lang)
		}
		if ip := api.MetaIp(ctx); ip != "" {
			h.logger.Debugf("client IP: %s", ip)
		}
	}

	return &pb.GetUserResponse{Id: req.Id, Name: "Alice", Email: "alice@example.com"}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)
	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
