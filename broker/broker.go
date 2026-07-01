package broker

import "github.com/starfork/stargo/util/pm"

type Broker interface {
	Publish(topic string, message Message) error
	Subscribe(topic string, handler MessageHandler) (Subscription, error)
	UnSubscribe() error
	Flush() error
}

type JetStreamBroker interface {
	Broker

	JetStreamPublish(topic string, msg Message) (*PubAck, error)
	JetStreamSubscribe(topic, consumerGroup string, handler JetStreamHandler) (Subscription, error)
}

type PubAck struct {
	Stream    string
	Sequence  uint64
	Duplicate bool
}

type Subscription interface {
	Unsubscribe() error
}

type Message struct {
	Topic     string
	Reply     string
	Header    pm.Pm
	Body      []byte
	MsgID     string
	Timestamp int64
}

type MessageHandler func(Message)

type JetStreamHandler func(msg Message) error

var brokerFactories = make(map[string]func(*Config) (Broker, error))

func Register(name string, factory func(*Config) (Broker, error)) {
	brokerFactories[name] = factory
}

func NewBroker(name string, conf *Config) (Broker, error) {
	if f, ok := brokerFactories[name]; ok {
		return f(conf)
	}
	return nil, nil
}
