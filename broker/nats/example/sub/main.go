package main

import (
	"fmt"

	"github.com/starfork/stargo/broker"
	nb "github.com/starfork/stargo/broker/nats"
)

func main() {
	b := nb.NewBroker(&broker.Config{
		App:  "example",
		Host: "nats://127.0.0.1:4222",
		Name: "nats",
	})
	b.Subscribe("example.*", func(m broker.Message) {

		fmt.Println(m.Header["stargo_borker_topic"])
	})

}
