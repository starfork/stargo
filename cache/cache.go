package cache

import (
	"context"
)

type Cache interface {
	// Get gets a cached value by key.
	Get(ctx context.Context, key string) (any, error)
	// Put stores a key-value pair into cache.
	Put(ctx context.Context, key string, value any) error
	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error
	Scan(ctx context.Context, key string, data any) error
}

type Marshaler interface {
	Marshal() (data []byte, err error)
}

type Unmarshaler interface {
	Unmarshal(data []byte) error
}
