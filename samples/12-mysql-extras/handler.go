package main

import (
	"context"
	"time"

	"github.com/starfork/stargo"
	mysqlscope "github.com/starfork/stargo/contrib/store/mysql"
	"github.com/starfork/stargo/contrib/store/mysql/uid"
	"github.com/starfork/stargo/logger"
	pb "github.com/starfork/stargo/samples/proto/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type repo struct {
	db  *gorm.DB
	uid *uid.UID // UID 生成器 / UID generator
}

type handler struct {
	repo   *repo
	logger logger.Logger
	pb.UnimplementedSampleServiceServer
}

func NewHandler(app *stargo.App) *handler {
	h := &handler{
		logger: app.Logger(),
	}

	// Config-first: 从 app 获取 MySQL 实例
	// Config-first: get MySQL instance from app
	if db := app.Store("mysql"); db != nil {
		gdb := db.Instance().(*gorm.DB)

		// UID 生成器：基于 MySQL 行锁的分布式自增 ID（业务表 uid，步长 100）
		// UID generator: MySQL row-lock based distributed auto-increment IDs
		// 表结构: CREATE TABLE uid (business_id VARCHAR(64) PRIMARY KEY, max_id BIGINT, step INT)
		uidGen, err := uid.New(gdb,
			uid.Table("uid"),
			uid.ID("user"),
			uid.Len(50),
			uid.Logger(h.logger),
		)
		if err != nil {
			h.logger.Errorf("uid.New: %v", err)
		}

		h.repo = &repo{
			db:  gdb,
			uid: uidGen,
		}
		h.logger.Infof("mysql store connected, uid=%v", uidGen != nil)
	}

	return h
}

// GetUser 查询单个用户 / Query a single user
func (h *handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	h.logger.Infof("GetUser: id=%d", req.Id)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

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

	return &pb.GetUserResponse{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}

// ListUser 使用 GORM scope 分页查询 / Uses GORM scope for pagination
// mysqlscope.Page(page, size) 自动处理 offset 和 limit，无需手写
// mysqlscope.Page(page, size) handles offset/limit automatically
func (h *handler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUser: page=%d size=%d", req.Page, req.PageSize)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	var users []*pb.GetUserResponse
	var total int64

	// Page scope: 避免手写 offset/limit，page<=0 自动修正为 1
	// Page scope: avoids manual offset/limit; page<=0 auto-corrects to 1
	db := h.repo.db.WithContext(ctx).Table("users")
	db.Count(&total)
	if err := db.Scopes(mysqlscope.Page(uint32(req.Page), uint32(req.PageSize))).Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListUserResponse{Users: users, Total: int32(total)}, nil
}

// CreateUser 创建用户，演示 UID 生成和时间戳处理
// CreateUser: demonstrates UID generation and timestamp handling
func (h *handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	h.logger.Infof("CreateUser: name=%s email=%s", req.Name, req.Email)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	// UID 生成器获取分布式自增 ID / Get distributed auto-increment ID from UID generator
	var userID int64
	if h.repo.uid != nil {
		id, err := h.repo.uid.Get()
		if err != nil {
			// 超时则回退到 DB 自增 / fallback to DB auto-increment on timeout
			h.logger.Warnf("uid.Get timeout, fallback: %v", err)
		} else {
			userID = int64(id)
		}
	}

	// store.Now() 返回 Asia/Shanghai 时区的格式化时间字符串
	// store.Now() returns a formatted time string in Asia/Shanghai timezone
	// store.TFORMAT = "2006-01-02T15:04:05+08:00"
	// store.TIME_LOCATION = "Asia/Shanghai"
	user := map[string]any{
		"name":  req.Name,
		"email": req.Email,
	}
	if userID > 0 {
		user["id"] = userID
	}
	user["created_at"] = time.Now().Format(time.RFC3339)

	result := h.repo.db.WithContext(ctx).Table("users").Create(user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	id := userID
	if id == 0 {
		id = result.RowsAffected
	}
	return &pb.CreateUserResponse{Id: id}, nil
}

// ListUserByTime 使用 Timezome scope 按时间范围查询 / Uses Timezome scope for time-range queries
// mysqlscope.Timezome(map[string]int64{"from": t1, "to": t2}, "created_at")
// 自动生成 BETWEEN 条件，支持仅 from、仅 to、from+to 三种模式
// Generates BETWEEN conditions; supports from-only, to-only, and from+to modes
func (h *handler) ListUserByTime(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	h.logger.Infof("ListUserByTime: page=%d size=%d", req.Page, req.PageSize)

	if h.repo == nil || h.repo.db == nil {
		return nil, status.Error(codes.Unavailable, "mysql not configured")
	}

	var users []*pb.GetUserResponse
	var total int64

	// 示例：查询最近 7 天的用户 / Example: query users from last 7 days
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	tz := map[string]int64{
		"from": weekAgo.Unix(),
		"to":   now.Unix(),
	}

	db := h.repo.db.WithContext(ctx).Table("users")
	db.Count(&total)
	if err := db.Scopes(
		mysqlscope.Timezome(tz, "created_at"),
		mysqlscope.Page(uint32(req.Page), uint32(req.PageSize)),
	).Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListUserResponse{Users: users, Total: int32(total)}, nil
}
