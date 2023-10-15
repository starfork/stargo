package store

import "github.com/starfork/stargo/config"

type Store interface {
	Connect(*config.Config)
	GetInstance() any
	Close() //关闭连接
}

// func New() Store {

// }
