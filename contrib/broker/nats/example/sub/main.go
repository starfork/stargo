package main

import (
	"fmt"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/contrib/broker/nats"
)

func main() {
	b := nats.NewBroker(&broker.Config{
		App:  "example",
		Host: "nats://127.0.0.1:4222",
		Name: "nats",
	})
	b.Subscribe("example.*", func(m broker.Message) {

		fmt.Println(m.Header["stargo_borker_topic"])
	})

}
