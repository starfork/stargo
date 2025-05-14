package store

type Store interface {
	Connect(conf ...*Config)
	//GetInstance(conf ...*config.Config) any
	Close() //关闭连接
}

var TIME_LOCATION = "Asia/Shanghai" //上海
var TFORMAT = "2006-01-02T15:04:05+08:00"
var TZ1K = true //js是否传回来的时间戳，是否1000
