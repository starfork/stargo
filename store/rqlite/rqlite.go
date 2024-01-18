package rqlite

import (
	"fmt"

	"github.com/rqlite/gorqlite"
	"github.com/starfork/stargo/store"
)

type Rqlite struct {
	conn *gorqlite.Connection
	c    *store.Config
}

func NewRqlite(config *store.Config) store.Store {
	// if config.Timezome != "" {
	// 	TIME_LOCATION = config.Timezome
	// }
	// if config.Timeformat != "" {
	// 	TFORMAT = config.Timeformat
	// }

	return &Rqlite{
		c: config,
	}
}

// Connect
func (e *Rqlite) Connect(confs ...*store.Config) {
	conn, _ := gorqlite.Open(fmt.Sprintf("%s:%s", e.c.Host, e.c.Port))
	e.conn = conn
}

func (e *Rqlite) GetInstance(conf ...*store.Config) *gorqlite.Connection {

	return e.conn
}

func (e *Rqlite) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Rqlite) Prefix(prefix string) string {

	return prefix
}
