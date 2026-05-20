package broker

import "github.com/starfork/stargo/pm"

type Broker interface {
	Publish(topic string, message Message) error
	Subscribe(topic string, handler MessageHandler)
	UnSubscribe() error
	Flush() error
}

type Message struct {
	Topic  string
	Reply  string
	Header pm.Pm
	Body   []byte
}

type MessageHandler func(Message)

var brokerFactories = make(map[string]func(*Config) Broker)

func Register(name string, factory func(*Config) Broker) {
	brokerFactories[name] = factory
}

func NewBroker(name string, conf *Config) Broker {
	if f, ok := brokerFactories[name]; ok {
		return f(conf)
	}
	return nil
}
