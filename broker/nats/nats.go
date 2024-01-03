package nats

import "github.com/starfork/stargo/broker"

type NatsBroker struct {
	c broker.Config
}

func NewNatsBrokder(c broker.Config) broker.Broker {
	return &NatsBroker{c}
}

func (e *NatsBroker) Public(broker.Message) error {
	return nil
}
func (e *NatsBroker) Subscribe()   {}
func (e *NatsBroker) UnSubscribe() {}
