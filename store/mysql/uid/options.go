package uid

import (
	"github.com/starfork/stargo/logger"
	"gorm.io/gorm"
)

type CheckFunc func(num ...uint32) uint32
type SetpFunc func(num uint32) uint32

type Options struct {
	table  string //表名称
	db     *gorm.DB
	logger logger.Logger
	id     string //业务id
	len    uint32

	fun  []CheckFunc
	setp SetpFunc
}

type Option func(o *Options)

func Logger(c logger.Logger) Option {
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
func Setp(c SetpFunc) Option {
	return func(o *Options) {
		o.setp = c
	}
}

func DefaultOptions() Options {
	o := Options{
		table: "uid",
		id:    "user_auth",
		len:   100, //step一般需要设置大于100，不然每次启动服务发会运行两次getFromDB
	}
	return o
}
