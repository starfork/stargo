package store

import "public/pkg/config"

type Store interface {
	Connect(*config.ServerConfig)
	GetInstance() any
	Close() //关闭连接
}

// func New() Store {

// }
