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

	ctx context.Context
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
		ctx:    context.Background(),
	}
	return s
}

// 添加任务
func (e *Redis) Push(t *task.Task) error {
	value := t.Marshal()
	interval := time.Now().Unix() + t.Delay
	member := redis.Z{
		Score:  float64(interval), //执行时间
		Member: t.Subkey(),
	}

	// Use pipeline to make operation atomic
	pipe := e.rdc.Pipeline()
	pipe.ZAdd(e.ctx, e.name, member)
	pipe.Set(e.ctx, e.name+"."+t.Subkey(), value, 0)
	_, err := pipe.Exec(e.ctx)
	return err
}

func (e *Redis) Pop(t *task.Task) error {
	processingKey := e.name + ".processing"
	
	// Use pipeline to remove from both main and processing sets atomically
	pipe := e.rdc.Pipeline()
	pipe.ZRem(e.ctx, e.name, t.Subkey())
	pipe.ZRem(e.ctx, processingKey, t.Subkey())
	pipe.Del(e.ctx, e.name+"."+t.Subkey())
	_, err := pipe.Exec(e.ctx)
	return err
}

// redis里面，有序集合新增，即可实现update
func (e *Redis) Update(t *task.Task) error {
	processingKey := e.name + ".processing"
	value := t.Marshal()
	interval := time.Now().Unix() + t.Delay
	member := redis.Z{
		Score:  float64(interval),
		Member: t.Subkey(),
	}

	// Use pipeline to remove from processing and add to main set atomically
	pipe := e.rdc.Pipeline()
	pipe.ZRem(e.ctx, processingKey, t.Subkey())
	pipe.ZAdd(e.ctx, e.name, member)
	pipe.Set(e.ctx, e.name+"."+t.Subkey(), value, 0)
	_, err := pipe.Exec(e.ctx)
	return err
}

// redis里面，有序集合新增，即可实现update
func (e *Redis) Clear(key string) error {
	rs := e.rdc.ZRange(e.ctx, e.name, 0, -1)
	if rs.Err() != nil {
		return rs.Err()
	}
	
	// Use pipeline to clear all tasks atomically
	pipe := e.rdc.Pipeline()
	for _, v := range rs.Val() {
		pipe.ZRem(e.ctx, e.name, v)
		pipe.Del(e.ctx, e.name+"."+v)
	}
	_, err := pipe.Exec(e.ctx)
	return err
}

func (e *Redis) FetchJob(step int64) ([]string, error) {
	now := time.Now().Unix()
	e_unix := strconv.FormatInt(now, 10)
	opt := &redis.ZRangeBy{
		Min: "-inf", // All expired tasks
		Max: e_unix, // Current time
		Count: step, // Limit number of tasks
	}

	rs := e.rdc.ZRangeByScore(e.ctx, e.name, opt)
	return rs.Val(), rs.Err()

}

func (e *Redis) ClaimJob(step int64) ([]string, error) {
	now := time.Now().Unix()
	processingKey := e.name + ".processing"
	
	// Lua script to atomically claim jobs
	luaScript := `
		local jobs = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
		if #jobs == 0 then
			return {}
		end
		
		local claimed = {}
		for i, job in ipairs(jobs) do
			-- Move job from main zset to processing zset with visibility timeout
			redis.call('ZREM', KEYS[1], job)
			redis.call('ZADD', KEYS[2], ARGV[3], job)
			table.insert(claimed, job)
		end
		
		return claimed
	`
	
	visibilityTimeout := int64(30) // 30 seconds visibility timeout
	args := []interface{}{
		strconv.FormatInt(now, 10),          // max score
		step,                                 // limit
		strconv.FormatInt(now+visibilityTimeout, 10), // new score for processing
	}
	
	result, err := e.rdc.Eval(e.ctx, luaScript, []string{e.name, processingKey}, args...).Result()
	if err != nil {
		return nil, err
	}
	
	// Convert result to string slice
	switch v := result.(type) {
	case []interface{}:
		res := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				res[i] = str
			}
		}
		return res, nil
	default:
		return nil, nil
	}
}
func (e *Redis) ReadTask(key string) (*task.Task, error) {
	rs := e.rdc.Get(e.ctx, e.name+"."+key)
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	task, err := task.UnmarshalJson(rs.Val())
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (e *Redis) ReclaimExpired() ([]string, error) {
	now := time.Now().Unix()
	processingKey := e.name + ".processing"
	
	// Lua script to reclaim expired jobs from processing back to main set
	luaScript := `
		local expired = redis.call('ZRANGEBYSCORE', KEYS[2], '-inf', ARGV[1])
		if #expired == 0 then
			return {}
		end
		
		local reclaimed = {}
		for i, job in ipairs(expired) do
			-- Get the payload
			local payload = redis.call('GET', KEYS[1] .. '.' .. job)
			if payload then
				-- Move back to main set with current time as score
				redis.call('ZADD', KEYS[1], ARGV[1], job)
				redis.call('ZREM', KEYS[2], job)
				table.insert(reclaimed, job)
			else
				-- No payload, just remove from processing
				redis.call('ZREM', KEYS[2], job)
			end
		end
		
		return reclaimed
	`
	
	result, err := e.rdc.Eval(e.ctx, luaScript, []string{e.name, processingKey}, now).Result()
	if err != nil {
		return nil, err
	}
	
	switch v := result.(type) {
	case []interface{}:
		res := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				res[i] = str
			}
		}
		return res, nil
	default:
		return nil, nil
	}
}
