package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Get gets a cached value by key.
	Get(ctx context.Context, key string) (*Item, error)
	// Put stores a key-value pair into cache.
	Put(ctx context.Context, key string, value *Item) error
	// Delete removes a key from cache.
	Delete(ctx context.Context, key string) error
}

type Item struct {
	Value      any
	Expiration int64
}

// Expired returns true if the item has expired.
func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

func (i *Item) String() string {
	return i.Value.(string)
}

func (i *Item) Float64() float64 {
	return i.Value.(float64)
}
