package stargo

import (
	"context"
	"time"

	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store"
	"go.uber.org/zap"

	sf "github.com/sony/sonyflake"
)

// 执行方法，按照规则制定 namespace.[service].XxxHandler
func (s *App) Invoke(ctx context.Context, app, method string, in, rs interface{}, h ...string) error {

	if s.client == nil {
		s.client = client.New(s.conf)
	}
	return s.client.Invoke(ctx, app, method, in, rs, h...)
}

func (s *App) GetLogger() *zap.SugaredLogger {
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

// func (s *App) GetMysql() *mysql.Mysql {
// 	if s.mysql == nil {
// 		s.mysql = mysql.Connect(s.conf)
// 		return s.mysql
// 	}
// 	return s.mysql
// }

// func (s *App) GetRedis() *redis.Redis {
// 	if s.redis == nil {
// 		s.redis = redis.Connect(s.conf)
// 		return s.redis
// 	}
// 	return s.redis

// }

// func (s *App) GetMongo() *mongo.Mongo {
// 	if s.mongo == nil {
// 		s.mongo = mongo.Connect(s.conf)
// 		return s.mongo
// 	}
// 	return s.mongo
// }

func (s *App) GetConfig() *config.Config {
	return s.conf
}

func (s *App) GetSfid(conf ...sf.Settings) *sf.Sonyflake {
	if s.sfid != nil {
		return s.sfid
	}
	st := sf.Settings{}
	if len(conf) > 0 {
		st = conf[0]
	}
	st.StartTime = time.Date(2021, 1, 18, 0, 0, 0, 0, time.UTC)
	return sf.NewSonyflake(st)
}
