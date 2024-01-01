package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/queue"
	"go.uber.org/zap"

	sredis "github.com/starfork/stargo/store/redis"
)

/**
通过redis的有续集和实现的一个任务延迟队列
*/

type Redis struct {
	rdc    *redis.Client
	name   string
	logger *zap.SugaredLogger
	ctx    context.Context
}

// 预留redis的配置
type RedisConfig struct {
}

func New(name string, app *stargo.App) queue.Store {

	s := &Redis{
		rdc:    app.Store("redis").(*sredis.Redis).GetInstance(),
		name:   name,
		logger: app.GetLogger(),
		ctx:    context.Background(),
	}

	return s
}

// 添加任务
func (e *Redis) Push(t *queue.Task) error {
	value := t.MarshalJson()
	interval := time.Now().Unix() + t.Delay
	member := redis.Z{
		Score:  float64(interval), //执行时间
		Member: t.Subkey(),
	}
	if rs := e.rdc.ZAdd(e.ctx, e.name, member); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Set(e.ctx, e.name+"."+t.Subkey(), value, 0); rs.Err() != nil {
		return rs.Err()
	}
	return nil

}

func (e *Redis) Pop(t *queue.Task) error {
	//e.logger.Debug("RemoveTask:", subkey)
	if rs := e.rdc.ZRem(e.ctx, e.name, t.Subkey()); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Del(e.ctx, e.name+"."+t.Subkey()); rs.Err() != nil {
		return rs.Err()
	}
	return nil
}

// redis里面，有序集合新增，即可实现update
func (e *Redis) Update(t *queue.Task) error {
	return e.Push(t)
}

func (e *Redis) FetchJob(step int64) ([]string, error) {
	now := time.Now().Unix()
	s_unix := strconv.FormatInt(now-step, 10)
	e_unix := strconv.FormatInt(now, 10)
	opt := &redis.ZRangeBy{
		Min: s_unix, //1秒前
		Max: e_unix, //当前时间
	}

	rs := e.rdc.ZRangeByScore(e.ctx, e.name, opt)
	return rs.Val(), rs.Err()

}
func (e *Redis) ReadTask(key string) (*queue.Task, error) {
	rs := e.rdc.Get(e.ctx, e.name+"."+key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	task, err := queue.UnmarshalJson(rs.Val())
	if err != nil {
		return nil, err
	}
	return task, nil
}
