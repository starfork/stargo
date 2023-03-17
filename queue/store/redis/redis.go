package storage

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Storage struct {
	rdc    *redis.Client
	name   string
	logger *zap.SugaredLogger
}

// 预留redis的配置
type RedisConfig struct {
}

func New(name string, rdc *redis.Client, logger ...*zap.SugaredLogger) *Storage {
	s := &Storage{rdc: rdc, name: name}
	if len(logger) > 0 {
		s.logger = logger[0]
	}
	return s
}

func (e *Storage) AddJob(name, key, value string, interval float64) error {
	subkey := name + "." + key
	member := redis.Z{
		Score:  interval, //执行时间
		Member: subkey,
	}
	e.rdc.ZAdd(e.name, member)
	e.rdc.Set(e.name+"."+subkey, value, 0)
	e.logger.Debugf("添加定时任务:%s,%f \r\n", subkey, interval)
	return nil
	//fmt.Println("添加定时任务", funcname)
}

func (e *Storage) FetchJob() []string {
	now := time.Now().Unix()
	s_unix := strconv.FormatInt(now-1, 10)
	e_unix := strconv.FormatInt(now, 10)
	opt := redis.ZRangeBy{
		Min: s_unix, //1秒前
		Max: e_unix, //当前时间
	}
	rs := e.rdc.ZRangeByScore(e.name, opt)
	return rs.Val()

}
func (e *Storage) FetchTask(name string) string {
	//fmt.Println("redis task:", name)
	return e.rdc.Get(e.name + "." + name).Val()
}

func (e *Storage) RemoveTask(name, key string) error {
	subkey := name + "." + key
	e.logger.Debug("RemoveTask:", subkey)
	e.rdc.ZRem(e.name, subkey)
	e.rdc.Del(e.name + "." + subkey)
	//fmt.Println(rs)
	return nil
}
