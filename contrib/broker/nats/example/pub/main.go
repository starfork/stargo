package main

import (
	"fmt"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/contrib/broker/nats"
	"github.com/starfork/stargo/pm"
)

func main() {
	b := nats.NewBroker(&broker.Config{
		App:  "example",
		Host: "127.0.0.1:4222",
		Name: "nats",
	})
	err := b.Publish("test", broker.Message{
		Header: pm.Pm{
			"abc": "12312",
		},
		Body: []byte("hello nats"),
	})
	b.Flush() //测试，重要
	fmt.Println(err)
}
