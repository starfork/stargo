package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/queue/task"
)

/**
通过redis的有续集和实现的一个任务延迟队列
*/

type Redis struct {
	rdc    *redis.Client
	name   string
	logger logger.Logger
}

// 预留redis的配置
type RedisConfig struct {
}

func New(rdc *redis.Client, opts ...store.Option) store.Store {
	options := store.DefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	s := &Redis{
		rdc:    rdc,
		name:   options.Name,
		logger: options.Logger,
	}
	return s
}

// 添加任务
func (e *Redis) Push(t *task.Task, pctx ...context.Context) error {
	value := t.MarshalJson()
	interval := time.Now().Unix() + t.Delay
	member := redis.Z{
		Score:  float64(interval), //执行时间
		Member: t.Subkey(),
	}
	var ctx context.Context
	if len(pctx) > 0 {
		ctx = pctx[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
	}

	rs := e.rdc.ZAdd(ctx, e.name, member)
	// fmt.Println(t.Subkey(), e.name, member, t.Delay)
	// fmt.Println("e.name", e.name, "---ZAdd result:", rs.Val(), "err:", rs.Err())
	if rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Set(ctx, e.name+"."+t.Subkey(), value, 0); rs.Err() != nil {
		return rs.Err()
	}
	return nil

}

func (e *Redis) Pop(t *task.Task, pctx ...context.Context) error {
	var ctx context.Context
	if len(pctx) > 0 {
		ctx = pctx[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
	}

	//e.logger.Debug("RemoveTask:", subkey)
	if rs := e.rdc.ZRem(ctx, e.name, t.Subkey()); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Del(ctx, e.name+"."+t.Subkey()); rs.Err() != nil {
		return rs.Err()
	}
	return nil
}

// redis里面，有序集合新增，即可实现update
func (e *Redis) Update(t *task.Task, pctx ...context.Context) error {
	return e.Push(t, pctx...)
}

// redis里面，有序集合新增，即可实现update
func (e *Redis) Clear(key string, pctx ...context.Context) error {
	var ctx context.Context
	if len(pctx) > 0 {
		ctx = pctx[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
	}

	rs := e.rdc.ZRange(ctx, e.name, 0, -1)
	if rs.Err() != nil {
		return rs.Err()
	}
	for _, v := range rs.Val() {
		e.Pop(&task.Task{
			Tag: v,
			Key: key,
		})
	}
	return nil
}

func (e *Redis) FetchJob(step int64, pctx ...context.Context) ([]string, error) {
	now := time.Now().Unix()
	s_unix := strconv.FormatInt(now-step, 10)
	e_unix := strconv.FormatInt(now, 10)
	opt := &redis.ZRangeBy{
		Min: s_unix, //1秒前
		Max: e_unix, //当前时间
	}

	var ctx context.Context
	if len(pctx) > 0 {
		ctx = pctx[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
	}

	rs := e.rdc.ZRangeByScore(ctx, e.name, opt)
	return rs.Val(), rs.Err()

}
func (e *Redis) ReadTask(key string, pctx ...context.Context) (*task.Task, error) {
	var ctx context.Context
	if len(pctx) > 0 {
		ctx = pctx[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
	}

	rs := e.rdc.Get(ctx, e.name+"."+key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	task, err := task.UnmarshalJson(rs.Val())
	if err != nil {
		return nil, err
	}
	return task, nil
}
