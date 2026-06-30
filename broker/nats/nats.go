package nats

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
)

type natsSubscription struct {
	sub *nats.Subscription
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
}

func init() {
	broker.Register("nats", func(c *broker.Config) (broker.Broker, error) {
		return NewBroker(c)
	})
}

func NewBroker(c *broker.Config) (broker.Broker, error) {
	// Add reconnected and disconnected handlers
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

	return &NatsBroker{c, nc}, nil
}

func (e *NatsBroker) Publish(topic string, msg broker.Message) error {
	b, err := jsoniter.Marshal(msg)
	if err != nil {
		return err
	}
	// Use consistent topic prefix
	return e.nc.Publish(e.c.App+"."+topic, b)
}

func (e *NatsBroker) Flush() error {
	return e.nc.Flush()
}

func (e *NatsBroker) Subscribe(topic string, handler broker.MessageHandler) (broker.Subscription, error) {
	// Use consistent topic prefix
	sub, err := e.nc.Subscribe(e.c.App+"."+topic, func(msg *nats.Msg) {
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
