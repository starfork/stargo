package storage

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/starfork/stargo/queue"
	"go.uber.org/zap"
)

/**
通过redis的有续集和实现的一个任务延迟队列
*/

type Redis struct {
	rdc    *redis.Client
	name   string
	logger *zap.SugaredLogger
}

// 预留redis的配置
type RedisConfig struct {
}

func New(name string, rdc *redis.Client, logger ...*zap.SugaredLogger) queue.Store {
	s := &Redis{rdc: rdc, name: name}
	if len(logger) > 0 {
		s.logger = logger[0]
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
	if rs := e.rdc.ZAdd(e.name, member); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Set(e.name+"."+t.Subkey(), value, 0); rs.Err() != nil {
		return rs.Err()
	}
	return nil

}

func (e *Redis) Pop(t *queue.Task) error {
	//e.logger.Debug("RemoveTask:", subkey)
	if rs := e.rdc.ZRem(e.name, t.Subkey()); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Del(e.name + "." + t.Subkey()); rs.Err() != nil {
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
	opt := redis.ZRangeBy{
		Min: s_unix, //1秒前
		Max: e_unix, //当前时间
	}
	rs := e.rdc.ZRangeByScore(e.name, opt)
	return rs.Val(), rs.Err()

}
func (e *Redis) ReadTask(key string) (*queue.Task, error) {
	//fmt.Println("redis task:", name)
	rs := e.rdc.Get(e.name + "." + key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	task, err := queue.UnmarshalJson(rs.Val())
	if err != nil {
		return nil, err
	}
	return task, nil
}
