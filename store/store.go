package store

import "github.com/starfork/stargo/config"

type Store interface {
	Connect(conf ...*config.Config)
	//GetInstance(conf ...*config.Config) any
	Close() //关闭连接
}

var TIME_LOCATION = "Asia/Shanghai" //上海
var TFORMAT = "2006-01-02T15:04:05+08:00"
