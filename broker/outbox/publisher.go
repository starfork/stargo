package outbox

import (
	"context"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
)

type Record struct {
	ID          int64  `json:"id"`
	AggregateID string `json:"aggregate_id"`
	EventType   string `json:"event_type"`
	Payload     []byte `json:"payload"`
	Topic       string `json:"topic"`
}

type Store interface {
	FetchPending(ctx context.Context, limit int) ([]*Record, error)
	MarkPublished(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

type OutboxPublisher struct {
	store    Store
	b        broker.Broker
	logger   logger.Logger
	interval time.Duration
}

func New(store Store, b broker.Broker) *OutboxPublisher {
	return &OutboxPublisher{
		store:    store,
		b:        b,
		interval: 1 * time.Second,
	}
}

func (p *OutboxPublisher) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.publishPending(ctx)
		}
	}
}

func (p *OutboxPublisher) publishPending(ctx context.Context) {
	records, err := p.store.FetchPending(ctx, 100)
	if err != nil {
		p.log("outbox fetch failed: %v", err)
		return
	}

	for _, r := range records {
		if err := p.b.Publish(r.Topic, broker.Message{
			Body: r.Payload,
		}); err != nil {
			p.log("outbox publish failed: %v, err: %v", r.ID, err)
			p.store.MarkFailed(ctx, r.ID)
		} else {
			p.store.MarkPublished(ctx, r.ID)
		}
	}
}

func (p *OutboxPublisher) log(format string, args ...interface{}) {
	if p.logger != nil {
		p.logger.Debugf(format, args...)
	}
}
