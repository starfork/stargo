package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/starfork/stargo/store"
)

type Influxdb struct {
	client influxdb2.Client
	c      *store.Config
}

func NewInfluxdb(config *store.Config) store.Store {
	return &Influxdb{
		c: config,
	}
}

func (e *Influxdb) Connect(conf ...*store.Config) {
	c := e.c

	e.client = influxdb2.NewClient(c.Host, c.Auth)

}

func (e *Influxdb) GetInstance(conf ...*store.Config) influxdb2.Client {
	if len(conf) > 0 {
		e.Connect(conf...)
		return e.client
	}
	if e.client == nil {
		e.Connect()
	}
	return e.client
}

func (e *Influxdb) Close() {
	e.client.Close()
}
