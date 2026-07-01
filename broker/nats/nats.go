package nats

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
)

type natsSubscription struct {
	sub  *nats.Subscription
	conn *nats.Conn
}

func (s *natsSubscription) Unsubscribe() error {
	if s.sub != nil {
		return s.sub.Unsubscribe()
	}
	return nil
}

type NatsBroker struct {
	c  *broker.Config
	nc *nats.Conn
	js nats.JetStreamContext
}

func init() {
	broker.Register("nats", func(c *broker.Config) (broker.Broker, error) {
		return NewBroker(c)
	})
}

func NewBroker(c *broker.Config) (broker.Broker, error) {
	opts := []nats.Option{
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.DefaultLogger.Infof("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			logger.DefaultLogger.Warnf("NATS disconnected from %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.DefaultLogger.Infof("NATS connection closed")
		}),
	}

	nc, err := nats.Connect(c.Host, opts...)
	if err != nil {
		return nil, err
	}

	b := &NatsBroker{c: c, nc: nc}

	if c.JetStream != nil && c.JetStream.Enabled {
		js, err := nc.JetStream()
		if err != nil {
			nc.Close()
			return nil, fmt.Errorf("jetstream init: %w", err)
		}
		b.js = js
		if err := b.ensureStream(); err != nil {
			nc.Close()
			return nil, fmt.Errorf("ensure stream: %w", err)
		}
		if c.JetStream.DLQEnabled {
			b.ensureDLQStream()
		}
	}

	return b, nil
}

func (e *NatsBroker) ensureStream() error {
	cfg := e.c.JetStream
	name := e.streamName()

	_, err := e.js.StreamInfo(name)
	if err == nil {
		_, err = e.js.UpdateStream(&nats.StreamConfig{
			Name:      name,
			Subjects:  cfg.Subjects,
			MaxAge:    cfg.MaxAge,
			MaxBytes:  cfg.MaxBytes,
			MaxMsgs:   cfg.MaxMsgs,
			Replicas:  cfg.Replicas,
		})
		return err
	}

	_, err = e.js.AddStream(&nats.StreamConfig{
		Name:      name,
		Subjects:  cfg.Subjects,
		MaxAge:    cfg.MaxAge,
		MaxBytes:  cfg.MaxBytes,
		MaxMsgs:   cfg.MaxMsgs,
		Replicas:  cfg.Replicas,
	})
	return err
}

func (e *NatsBroker) ensureDLQStream() {
	name := e.dlqStreamName()
	e.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{name + ".>"},
	})
}

func (e *NatsBroker) streamName() string {
	if e.c.JetStream != nil && e.c.JetStream.StreamName != "" {
		return e.c.JetStream.StreamName
	}
	return "stargo_" + e.c.App
}

func (e *NatsBroker) dlqStreamName() string {
	if e.c.JetStream != nil && e.c.JetStream.DLQName != "" {
		return e.c.JetStream.DLQName
	}
	return e.streamName() + "_dlq"
}

func (e *NatsBroker) fullTopic(topic string) string {
	return e.c.App + "." + topic
}

func (e *NatsBroker) Publish(topic string, msg broker.Message) error {
	b, err := jsoniter.Marshal(msg)
	if err != nil {
		return err
	}
	if e.js != nil {
		jsOpts := []nats.PubOpt{}
		if msg.MsgID != "" {
			jsOpts = append(jsOpts, nats.MsgId(msg.MsgID))
		}
		_, err = e.js.Publish(e.fullTopic(topic), b, jsOpts...)
		return err
	}
	return e.nc.Publish(e.fullTopic(topic), b)
}

func (e *NatsBroker) JetStreamPublish(topic string, msg broker.Message) (*broker.PubAck, error) {
	if e.js == nil {
		return nil, fmt.Errorf("jetstream not enabled")
	}
	b, err := jsoniter.Marshal(msg)
	if err != nil {
		return nil, err
	}
	jsOpts := []nats.PubOpt{}
	if msg.MsgID != "" {
		jsOpts = append(jsOpts, nats.MsgId(msg.MsgID))
	}
	ack, err := e.js.Publish(e.fullTopic(topic), b, jsOpts...)
	if err != nil {
		return nil, err
	}
	return &broker.PubAck{
		Stream:    ack.Stream,
		Sequence:  ack.Sequence,
		Duplicate: ack.Duplicate,
	}, nil
}

func (e *NatsBroker) JetStreamSubscribe(topic, consumerGroup string, handler broker.JetStreamHandler) (broker.Subscription, error) {
	if e.js == nil {
		return nil, fmt.Errorf("jetstream not enabled")
	}

	cfg := e.c.JetStream
	consumerName := consumerGroup + "_" + topic
	ackWait := cfg.AckWait
	if ackWait <= 0 {
		ackWait = 30 * time.Second
	}
	maxDeliver := cfg.MaxDeliver
	if maxDeliver <= 0 {
		maxDeliver = 7
	}

	sub, err := e.js.QueueSubscribe(e.fullTopic(topic), consumerGroup, func(msg *nats.Msg) {
		bmsg := broker.Message{}
		if err := jsoniter.Unmarshal(msg.Data, &bmsg); err != nil {
			msg.Term()
			return
		}
		bmsg.Topic = msg.Subject

		if err := handler(bmsg); err != nil {
			meta, _ := msg.Metadata()
			if meta != nil && meta.NumDelivered >= uint64(maxDeliver) {
				if cfg.DLQEnabled {
					dlqTopic := e.dlqStreamName() + "." + topic
					e.js.Publish(dlqTopic, msg.Data)
				}
				msg.Term()
			} else {
				msg.Nak()
			}
		} else {
			msg.Ack()
		}
	}, nats.ManualAck(), nats.AckWait(ackWait), nats.MaxDeliver(maxDeliver), nats.Durable(consumerName))

	if err != nil {
		return nil, err
	}

	return &natsSubscription{sub: sub, conn: e.nc}, nil
}

func (e *NatsBroker) Flush() error {
	return e.nc.Flush()
}

func (e *NatsBroker) Subscribe(topic string, handler broker.MessageHandler) (broker.Subscription, error) {
	sub, err := e.nc.Subscribe(e.fullTopic(topic), func(msg *nats.Msg) {
		bmsg := broker.Message{}
		err := jsoniter.Unmarshal(msg.Data, &bmsg)
		if err != nil {
			logger.DefaultLogger.Errorf("NATS unmarshal error: %v", err)
			return
		}
		bmsg.Topic = msg.Subject
		handler(bmsg)
	})
	if err != nil {
		return nil, err
	}

	return &natsSubscription{sub: sub}, nil
}

func (e *NatsBroker) UnSubscribe() error {
	return e.nc.Drain()
}
