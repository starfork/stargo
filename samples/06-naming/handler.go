package main

import (
	"context"
	"fmt"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
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

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	// Access the registry to discover other services.
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

	// The Service() method returns the current service's registration info.
	svc := h.app.Service()
	h.logger.Infof("service registered: %s @ %s", svc.Name, svc.Addr)

	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}

func (h *handler) resolveTarget() string {
	if r := h.app.Resolver(); r != nil {
		return fmt.Sprintf("%s:///%s/%s", r.Scheme(), r.Config().Org, "target-service")
	}
	return ""
}
