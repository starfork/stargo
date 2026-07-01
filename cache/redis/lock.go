package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starfork/stargo/cache"
)

var ErrLockHeld = errors.New("lock already held")

type RedisLocker struct {
	client *redis.Client
}

func NewLocker(client *redis.Client) cache.Locker {
	return &RedisLocker{client: client}
}

type redisUnlocker struct {
	client *redis.Client
	key    string
	token  string
}

func randomToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (l *RedisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (cache.Unlocker, error) {
	token := randomToken()
	lockKey := "lock:" + key

	for {
		ok, err := l.client.SetNX(ctx, lockKey, token, ttl).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			return &redisUnlocker{client: l.client, key: lockKey, token: token}, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (l *RedisLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (cache.Unlocker, error) {
	token := randomToken()
	lockKey := "lock:" + key

	ok, err := l.client.SetNX(ctx, lockKey, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLockHeld
	}
	return &redisUnlocker{client: l.client, key: lockKey, token: token}, nil
}

func (u *redisUnlocker) Unlock(ctx context.Context) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	_, err := u.client.Eval(ctx, script, []string{u.key}, u.token).Result()
	return err
}
