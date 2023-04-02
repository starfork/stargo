package stargo

import (
	"context"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store/mysql"
	"github.com/starfork/stargo/store/redis"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 执行方法，按照规则制定 namespace.[service].XxxHandler
func (s *App) Invoke(ctx context.Context, app, method string, in, rs interface{}) error {
	conf := s.conf
	r := naming.NewResolver(conf.Registry)

	//统一独立部署，只有一个target
	target := app
	if s.opts.Config.Deploy == config.DEPLOY_Monolithic {
		target = s.opts.Name
	}

	conn, err := grpc.Dial(s.registry.Scheme()+"://"+target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(r))
	if err != nil {
		return err
	}

	handler := cases.Title(language.English).String(app) + "Handler"

	return conn.Invoke(ctx, "/"+s.opts.Org+"."+app+"."+handler+"/"+method, in, rs)

	//return s.invoke(ctx, app, method, cases.Title(language.English).String(app)+"Handler", in, rs)
}

func (s *App) GetLogger() *zap.SugaredLogger {
	return s.logger
}

func (s *App) GetMysql() *mysql.Mysql {
	if s.mysql == nil {
		return mysql.Connect(s.conf)
	}
	return s.mysql
}

func (s *App) GetRedis() *redis.Redis {
	if s.redis == nil {
		return redis.Connect(s.conf)
	}
	return s.redis
}

func (s *App) GetConfig() *config.Config {

	return s.config
}
