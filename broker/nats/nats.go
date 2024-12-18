package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/starfork/stargo/broker"
)

type NatsBroker struct {
	c  *broker.Config
	nc *nats.Conn
	js jetstream.JetStream
}

func NewBroker(c *broker.Config) broker.Broker {
	nc, err := nats.Connect(c.Host)
	if err != nil {
		panic(err.Error())
	}
	js, _ := jetstream.New(nc)

	return &NatsBroker{c, nc, js}
}

func (e *NatsBroker) Public(broker.Message) error {
	return nil
}
func (e *NatsBroker) Subscribe()   {}
func (e *NatsBroker) UnSubscribe() {}
