package config

import (
	"strings"
)

type Val map[string]string

type Config struct {
	val Val

	//Option Options
	store StoreInterface
}

type KV struct {
	Key string
	Val string
}

type StoreInterface interface {
	Load() []*KV
	Set(pfx string, value map[string]string) error
	//Get(name string) string
}

func New(opts ...Option) *Config {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	c := &Config{
		val:   make(Val),
		store: options.Store,
	}
	c.Load()
	return c
}

// 从store里面夹在，然后放到val里面
func (e *Config) Load() {
	result := e.store.Load()
	for _, v := range result {
		e.val[v.Key] = v.Val
	}
}

func (e *Config) Set(pfx string, value map[string]string) error {

	for k, v := range value {
		key := strings.ToUpper(pfx) + "_" + strings.ToUpper(k)
		e.SetVal(key, v)
	}

	return e.store.Set(pfx, value)
}

// for tet
func (e *Config) SetVal(key string, value string) {
	if e.val == nil {
		e.val = make(Val)
	}
	e.val[key] = value
}
