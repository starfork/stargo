package uid

import "gorm.io/gorm"

// 日志接口
type Log interface {
	Error(erargs ...interface{})
}

type CheckFunc func(num ...uint32) uint32

type Options struct {
	table  string //表名称
	db     *gorm.DB
	logger Log
	id     string //业务id
	len    uint32

	fun []CheckFunc
}

type Option func(o *Options)

func Logger(c Log) Option {
	return func(o *Options) {
		o.logger = c
	}
}

func DB(c *gorm.DB) Option {
	return func(o *Options) {
		o.db = c
	}
}
func Len(c uint32) Option {
	return func(o *Options) {
		o.len = c
	}
}
func ID(c string) Option {
	return func(o *Options) {
		o.id = c
	}
}
func Table(c string) Option {
	return func(o *Options) {
		o.table = c
	}
}

// 可多次调用
func Fun(c CheckFunc) Option {
	return func(o *Options) {
		o.fun = append(o.fun, c)
	}
}

func DefaultOptions() Options {
	o := Options{
		table: "uid",
		id:    "user_auth",
		len:   100,
	}
	return o
}
