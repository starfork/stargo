package main

import (
	"context"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/util/pm"
	"github.com/starfork/stargo/queue"
	qstore "github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/queue/task"
	qredis "github.com/starfork/stargo/queue/store/redis"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"github.com/redis/go-redis/v9"
)

type handler struct {
	app    *stargo.App
	logger logger.Logger
	queue  *queue.Queue
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	h := &handler{
		app:    app,
		logger: app.Logger(),
	}

	// 从 Redis store 初始化队列存储 / Init queue store from Redis
	if rdc := app.Store("redis"); rdc != nil {
		client := rdc.Instance().(*redis.Client)
		store := qredis.New(
			client,
			qstore.WithName("demo-queue"),
		)

		// 创建队列引擎 / Create queue engine
		q := queue.New(
			store,
			queue.WithStep(1),             // 轮询跨度 / Poll step
			queue.WithInterval(1),         // 轮询间隔(秒) / Poll interval (seconds)
			queue.WithMaxTrhead(5),        // 最大并发 / Max concurrency
		)

		// 注册任务处理器 / Register task handler
		q.Register("send_email", func(tk *task.Task) error {
			h.logger.Infof("processing task: key=%s tag=%s", tk.Key, tk.Tag)
			// 实际业务逻辑 / Actual business logic
			return nil
		})

		h.queue = q
	}

	return h
}

func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)
	return &pb.GetUserResponse{Id: req.Id, Name: "Alice", Email: "alice@example.com"}, nil
}

// CreateUser 创建用户后推送延迟任务 / Creates user and pushes a delayed task
func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	if h.queue != nil {
		// 5 秒后发送欢迎邮件 / Send welcome email after 5 seconds
		tk := &task.Task{
			Key:      "user_welcome",
			Tag:      "send_email",
			User:     req.Email,
			Delay:    5,                          // 延迟秒数 / Delay in seconds
			DelayTTL: []int64{3600},              // 任务最长存活 / Max task TTL
			RetryMax: 3,                          // 最大重试 / Max retries
			Args: pm.Pm{
				"name":  req.Name,
				"email": req.Email,
			},
		}
		if err := h.queue.Push(tk); err != nil {
			h.logger.Errorf("queue push error: %v", err)
		} else {
			h.logger.Infof("queued welcome email for %s", req.Email)
		}
	}

	return &pb.CreateUserResponse{Id: 100}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)
	return &pb.ListUserResponse{Total: 0}, nil
}
