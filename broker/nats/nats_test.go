package nats

import (
	"fmt"
	"testing"

	"github.com/starfork/stargo/broker"
)

func TestPublish(t *testing.T) {
	b, err := NewBroker(&broker.Config{
		Host: "127.0.0.1:4222",
		Name: "nats",
	})
	if err != nil {
		t.Skipf("nats not available: %v", err)
		return
	}
	err = b.Publish("nats-test", broker.Message{
		Body: []byte("hello nats"),
	})
	fmt.Println(err)
}
