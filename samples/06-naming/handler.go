package main

import (
	"context"
	"fmt"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
)

// 服务启动时自动注册到 etcd, 停止时自动摘除 / Service auto-registers to etcd at Run() and auto-deregisters at Stop().

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

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	// 通过注册中心获取其他服务的发现信息 / Access the registry to discover other services.
	if reg := h.app.Registry(); reg != nil {
		h.logger.Debugf("registry scheme: %s", reg.Scheme())
	}

	return &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	// Service() 返回当前服务的注册信息 (名称+地址) / Service() returns the current service's registration info (name + address).
	svc := h.app.Service()
	h.logger.Infof("service registered: %s @ %s", svc.Name, svc.Addr)

	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}

// Resolver 将服务名解析为 gRPC 目标地址格式 "scheme:///org/service"
// Resolver translates a service name to the gRPC target URI format "scheme:///org/service".
func (h *handler) resolveTarget() string {
	if r := h.app.Resolver(); r != nil {
		return fmt.Sprintf("%s:///%s/%s", r.Scheme(), r.Config().Org, "target-service")
	}
	return ""
}
