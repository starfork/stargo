package main

import (
	"context"
	"fmt"
	"time"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	app    *stargo.App
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	h := &handler{
		app:    app,
		logger: app.Logger(),
	}

	// Subscribe to an event on startup.
	if b := app.Broker(); b != nil {
		b.Subscribe("user.created", func(msg broker.Message) {
			h.logger.Infof("event received: topic=%s body=%s", msg.Topic, string(msg.Body))
		})
		h.logger.Infof("subscribed to user.created")
	}

	return h
}

// CreateUser publishes a broker event after creating the user.
func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	id := int64(time.Now().UnixNano())

	// Publish an event via the broker.
	if b := h.app.Broker(); b != nil {
		msg := broker.Message{
			Topic: "user.created",
			Body:  []byte(fmt.Sprintf(`{"id":%d,"name":"%s"}`, id, req.Name)),
		}
		if err := b.Publish("user.created", msg); err != nil {
			h.logger.Errorf("publish user.created: %v", err)
		}
	}

	return &pb.CreateUserResponse{Id: id}, nil
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)
	return &pb.GetUserResponse{Id: req.Id, Name: "Alice", Email: "alice@example.com"}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
