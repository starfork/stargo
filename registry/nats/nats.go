package nats

import (
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/starfork/micro-boot/registry"
)

type natsRegistry struct {
	addrs      []string
	opts       registry.Options
	nopts      nats.Options
	queryTopic string
	watchTopic string

	sync.RWMutex
	conn      *nats.Conn
	services  map[string][]*registry.Service
	listeners map[string]chan bool
}
