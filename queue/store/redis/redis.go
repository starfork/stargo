package storage

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/starfork/stargo/queue"
	"go.uber.org/zap"
)

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
	name := e.name + "_" + t.Tag
	value := t.MarshalJson()
	interval := time.Now().Unix() + t.Delay
	subkey := t.Tag + "." + t.Key
	member := redis.Z{
		Score:  float64(interval), //执行时间
		Member: subkey,
	}
	if rs := e.rdc.ZAdd(name, member); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Set(e.name+"."+subkey, value, 0); rs.Err() != nil {
		return rs.Err()
	}
	return nil

}

func (e *Redis) Pop(t *queue.Task) error {
	name := e.name + "_" + t.Tag
	subkey := t.Tag + "." + t.Key
	//e.logger.Debug("RemoveTask:", subkey)
	if rs := e.rdc.ZRem(name, subkey); rs.Err() != nil {
		return rs.Err()
	}
	if rs := e.rdc.Del(name + "." + subkey); rs.Err() != nil {
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
