package broker

import (
	"context"
	"time"
)

type IdempotentStore interface {
	IsProcessed(ctx context.Context, msgID string) (bool, error)
	MarkProcessed(ctx context.Context, msgID string, ttl time.Duration) error
}

func WithIdempotent(store IdempotentStore, handler JetStreamHandler) JetStreamHandler {
	return func(msg Message) error {
		if msg.MsgID != "" {
			processed, _ := store.IsProcessed(context.Background(), msg.MsgID)
			if processed {
				return nil
			}
		}
		err := handler(msg)
		if err == nil && msg.MsgID != "" {
			store.MarkProcessed(context.Background(), msg.MsgID, 24*time.Hour)
		}
		return err
	}
}
