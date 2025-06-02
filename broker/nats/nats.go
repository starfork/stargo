package nats

import (
	jsoniter "github.com/json-iterator/go"
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
	b, err := jsoniter.Marshal(msg)
	if err != nil {
		return err
	}
	return e.nc.Publish(e.c.App+"."+topic, b)
	//
}

func (e *NatsBroker) Flush() error {
	return e.nc.Flush()
}

// 需要完整的app.name.
// !大多数情况下，需要使用 go Subscribe
func (e *NatsBroker) Subscribe(topic string, handler broker.MessageHandler) {

	e.nc.Subscribe(topic, func(msg *nats.Msg) {
		bmsg := broker.Message{}
		err := jsoniter.Unmarshal(msg.Data, &bmsg)
		if err == nil {
			handler(bmsg)
		}
	})

	select {}

}
func (e *NatsBroker) UnSubscribe() error {
	return e.nc.Drain()
}
