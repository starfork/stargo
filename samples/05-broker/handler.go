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

	// 通过 broker 订阅主题 / Subscribe to a topic via the message broker.
	if b := app.Broker(); b != nil {
		// Subscribe 注册一个主题监听器, broker.Message 包含 Topic 和 Body 字段
		// Subscribe registers a topic listener; broker.Message carries Topic and Body fields.
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

	// 通过 broker 发布事件 / Publish an event via the message broker.
	if b := h.app.Broker(); b != nil {
		// broker.Message 是标准的消息结构: Topic 主题 + Body 负载
		// broker.Message is the standard message envelope: Topic routing key + Body payload.
		msg := broker.Message{
			Topic: "user.created",
			Body:  []byte(fmt.Sprintf(`{"id":%d,"name":"%s"}`, id, req.Name)),
		}
		// Publish 将消息发送到 NATS 主题, 订阅方会异步收到 / Publish sends the message to the NATS subject; subscribers receive it asynchronously.
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
