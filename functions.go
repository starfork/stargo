package stargo

import (
	"context"

	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/store"
)

// 执行方法，按照规则制定 namespace.[service].XxxHandler
func (s *App) Invoke(ctx context.Context, app, method string, in, rs interface{}, h ...string) error {

	if s.client == nil {
		s.client = client.New(s.conf)
	}
	return s.client.Invoke(ctx, app, method, in, rs, h...)
}

func (s *App) GetLogger() logger.Logger {
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

func (s *App) GetConfig() *config.Config {
	return s.conf
}
