package cache

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (Unlocker, error)
	TryLock(ctx context.Context, key string, ttl time.Duration) (Unlocker, error)
}

type Unlocker interface {
	Unlock(ctx context.Context) error
}
