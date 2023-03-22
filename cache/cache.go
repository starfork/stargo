package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Get gets a cached value by key.
	Get(ctx context.Context, key string) (interface{}, time.Time, error)
	// Put stores a key-value pair into cache.
	Put(ctx context.Context, key string, val interface{}, d time.Duration) error
	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error
	// String returns the name of the implementation.
	String() string
}

type Item struct {
	Value      interface{}
	Expiration int64
}

// Expired returns true if the item has expired.
func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}

	return time.Now().UnixNano() > i.Expiration
}
