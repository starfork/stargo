package influxdb

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/starfork/stargo/store"
)

type Influxdb struct {
	client *influxdb3.Client
	c      *store.Config
}

func NewInfluxdb(config *store.Config) store.Store {
	return &Influxdb{
		c: config,
	}
}

func (e *Influxdb) Connect(conf ...*store.Config) {
	c := e.c
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     c.Host,
		Token:    c.Auth,
		Database: c.Name,
	})
	if err != nil {
		panic(err)
	}
	// Close client at the end and escalate error if present
	// defer func(client *influxdb3.Client) {
	// 	err := client.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }(client)
	e.client = client

}

func (e *Influxdb) GetInstance(conf ...*store.Config) *influxdb3.Client {
	if len(conf) > 0 {
		e.Connect(conf...)
		return e.client
	}
	return e.client
}

func (e *Influxdb) Close() {
	e.client.Close()
}
