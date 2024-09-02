package stargo

import (
	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
	"google.golang.org/grpc"
)

// // Server
func (s *App) RpcServer() *grpc.Server {
	return s.server.Server()
}

func (s *App) Client() *client.Client {
	if s.resolver != nil {
		return client.New(s.ctx, s.resolver, s.logger)
	}
	return nil
}

// // 返回标准服务格式
func (s *App) Service() naming.Service {
	return s.server.Service()
}

func (s *App) Registry() naming.Registry {
	return s.registry
}
func (s *App) Resolver() naming.Resolver {
	return s.resolver
}

func (s *App) Logger(l ...logger.Logger) logger.Logger {
	if len(l) > 0 {
		s.logger = l[0]
	}
	return s.logger
}

// // 获取或者创建一个store
func (s *App) Store(name string, st ...store.Store) store.Store {
	if len(st) > 0 {
		s.store[name] = st[0]
		return s.store[name]
	} else {
		if store, ok := s.store[name]; ok {
			return store
		}
	}
	return nil
}

// func (s *App) Mysql() store.Store {
// 	return s.Store("mysql")
// }
// func (s *App) Redis() store.Store {
// 	return s.Store("redis")
// }

// func (s *App) Config() *config.Config {
// 	return s.conf
// }

// func (s *App) Broker() broker.Broker {
// 	return s.broker
// }
