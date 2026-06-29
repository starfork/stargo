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

	// tracer.DefaultTracer 初始为 noop; 通过配置 tracer.driver 并 import jaeger 插件启用
	// tracer.DefaultTracer is a noop by default; enable via config tracer.driver + import jaeger plugin
	// 配置方式 / Configuration:
	//   config.yaml: tracer: { driver: jaeger, host: "127.0.0.1:6831", name: "trace-demo" }
	//   main.go:     import _ "github.com/starfork/stargo/tracer/jaeger"
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
