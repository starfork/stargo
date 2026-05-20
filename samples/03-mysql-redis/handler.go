package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type repo struct {
	db  *gorm.DB
	rdb *redis.Client
}

type handler struct {
	repo *repo
	log  logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	h := &handler{
		log: app.Logger(),
	}

	// Config-first: if YAML has store.mysql, it's auto-connected.
	if db := app.Store("mysql"); db != nil {
		h.repo = &repo{
			db: db.Instance().(*gorm.DB),
		}
		h.log.Infof("mysql store connected")
	}

	// If YAML has store.redis, it's auto-connected as well.
	if rdc := app.Store("redis"); rdc != nil {
		if h.repo == nil {
			h.repo = &repo{}
		}
		h.repo.rdb = rdc.Instance().(*redis.Client)
		h.log.Infof("redis store connected")
	}

	return h
}

// GetUser queries MySQL for a user, with Redis cache-aside.
func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.log.Infof("GetUser: id=%d", req.Id)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	// Cache-aside: check Redis first.
	if h.repo.rdb != nil {
		val, err := h.repo.rdb.Get(ctx, fmt.Sprintf("user:%d", req.Id)).Result()
		if err == nil {
			h.log.Debugf("cache hit for user:%d: %s", req.Id, val)
			return &pb.GetUserResponse{Id: req.Id, Name: val}, nil
		}
		h.log.Debugf("cache miss for user:%d", req.Id)
	}

	// Query MySQL.
	var user struct {
		ID    int64
		Name  string
		Email string
	}
	if err := h.repo.db.WithContext(ctx).
		Table("users").
		Where("id = ?", req.Id).
		First(&user).Error; err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Cache the result in Redis.
	if h.repo.rdb != nil {
		h.repo.rdb.Set(ctx, fmt.Sprintf("user:%d", req.Id), user.Name, 0)
	}

	return &pb.GetUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.log.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	result := h.repo.db.WithContext(ctx).
		Table("users").
		Create(map[string]any{
			"name":  req.Name,
			"email": req.Email,
		})
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &pb.CreateUserResponse{Id: result.RowsAffected}, nil
}

func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.log.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	var users []*pb.GetUserResponse
	offset := int((req.Page - 1) * req.PageSize)
	if err := h.repo.db.WithContext(ctx).
		Table("users").
		Offset(offset).
		Limit(int(req.PageSize)).
		Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var total int64
	h.repo.db.Table("users").Count(&total)

	return &pb.ListUserResponse{
		Users: users,
		Total: int32(total),
	}, nil
}
