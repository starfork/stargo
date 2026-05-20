package main

import (
	"time"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/pm"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("broker-demo", conf)

	// Broker is configured via YAML under the "broker" section.
	// The NATS implementation is used when the config is present.
	// Topic is prefixed with the app name (e.g., "broker-demo.demo.event").

	b := app.Broker()
	if b == nil {
		app.LogInfof("broker not configured, skipping demo")
		return
	}

	// Subscribe before publishing to receive the message
	b.Subscribe("demo.event", func(m broker.Message) {
		app.LogInfof("received: topic=%s body=%s", m.Topic, string(m.Body))
	})

	// Wait for subscription to be active
	time.Sleep(100 * time.Millisecond)

	// Publish a message
	msg := broker.Message{
		Topic:  "demo.event",
		Header: pm.Pm{"source": "broker-demo"},
		Body:   []byte(`{"hello":"world"}`),
	}
	if err := b.Publish("demo.event", msg); err != nil {
		app.LogErrorf("publish: %v", err)
	} else {
		app.LogInfof("published to demo.event")
	}

	time.Sleep(time.Second)
}
