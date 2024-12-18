package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Get gets a cached value by key.
	Get(ctx context.Context, key string) (any, error)
	// Fetch is batch version of get
	Fetch(ctx context.Context, key []string) ([]any, error)
	// Put stores a key-value pair into cache.
	Put(ctx context.Context, key string, value any, timeout ...time.Duration) error
	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error

	IsExist(ctx context.Context, key string) (bool, error)
	ClearAll(ctx context.Context) error

	Incr(ctx context.Context, key string) error
	// Decrement a cached int value by key, as a counter.
	Decr(ctx context.Context, key string) error

	//Expire(ctx context.Context, key string) error
}

type Marshaler interface {
	Marshal() (data []byte, err error)
}

type Unmarshaler interface {
	Unmarshal(data []byte) error
}
