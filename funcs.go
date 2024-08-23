package stargo

import (
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
	"google.golang.org/grpc"
)

// Server
func (s *App) RpcServer() *grpc.Server {
	return s.rpcServer
}

// Server
func (s *App) HttpServer() *grpc.Server {
	return s.rpcServer
}
// 返回标准服务格式
func (s *App) Service() naming.Service {
	return naming.Service{
		Org:  s.opts.Org,
		Name: s.opts.Name,
		Addr: s.conf.RpcServer.Host,
	}
}

func (s *App) Registry() naming.Registry {
	return s.registry
}

func (s *App) Logger(conf ...logger.Config) logger.Logger {
	return s.logger
}

// 获取或者创建一个store
func (s *App) Store(name string, st ...store.Store) store.Store {
	if len(st) > 0 {
		sto := st[0]
		s.store[name] = sto
		return sto
	} else {
		if store, ok := s.store[name]; ok {
			return store
		}
	}
	return nil
}

func (s *App) Mysql() store.Store {
	return s.Store("mysql")
}
func (s *App) Redis() store.Store {
	return s.Store("redis")
}

func (s *App) Config() *config.Config {
	return s.conf
}

func (s *App) Broker() broker.Broker {
	return s.broker
}
