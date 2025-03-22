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

func (e *NatsBroker) Publish(topic string, msg broker.Message) error {
	//必须带上当前的app名字
	return e.nc.Publish(e.c.App+"."+topic, msg.Body)
}
func (e *NatsBroker) Subscribe(topic string, handler broker.MessageHandler) {
	e.nc.Subscribe(topic, func(msg *nats.Msg) {
		handler(broker.Message{Body: msg.Data})
	})
}
func (e *NatsBroker) UnSubscribe() {}
