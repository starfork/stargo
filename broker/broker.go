package broker

type Broker interface {
	Publish(topic string, message Message) error
	Subscribe(topic string, handler MessageHandler)
	UnSubscribe()
}

type Message struct {
	Header map[string]any
	Body   []byte
}

type MessageHandler func(Message)
