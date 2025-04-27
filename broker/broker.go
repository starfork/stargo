package broker

import "github.com/starfork/stargo/pm"

type Broker interface {
	Publish(topic string, message Message) error
	Subscribe(topic string, handler MessageHandler)
	UnSubscribe() error
	Flush() error //短期测试（tester等场景）
}

type Message struct {
	Header pm.Pm
	Body   []byte
}

type MessageHandler func(Message)
