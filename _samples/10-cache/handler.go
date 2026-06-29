package main

import (
	"context"
	"fmt"
	"time"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/cache"
	redis_cache "github.com/starfork/stargo/cache/redis"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	logger logger.Logger
	cache  cache.Cache
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	h := &handler{
		logger: app.Logger(),
	}
	if rdc := app.Store("redis"); rdc != nil {
		h.cache = redis_cache.New(rdc)
		h.logger.Infof("redis cache ready")
	}
	return h
}

// GetUser 演示缓存使用 / Demonstrates cache-aside pattern
func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	// 检查缓存 / Check cache first
	key := fmt.Sprintf("user:%d", req.Id)
	if h.cache != nil {
		if val, err := h.cache.Get(ctx, key); err == nil {
			if name, ok := val.(string); ok {
				h.logger.Debugf("cache hit: %s", key)
				return &pb.GetUserResponse{Id: req.Id, Name: name}, nil
			}
		}
		h.logger.Debugf("cache miss: %s", key)
	}

	// 模拟数据库查询 / Simulate DB query
	resp := &pb.GetUserResponse{
		Id:    req.Id,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	// 回填缓存 / Backfill cache
	if h.cache != nil {
		if err := h.cache.Put(ctx, key, resp.Name, 5*time.Minute); err != nil {
			h.logger.Warnf("cache put error: %v", err)
		}
		h.logger.Debugf("cache set: %s", key)
	}

	return resp, nil
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
