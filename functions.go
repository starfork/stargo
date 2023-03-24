package stargo

import (
	"context"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store/mysql"
	"github.com/starfork/stargo/store/redis"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
)

// 执行方法，按照规则制定 namespace.[service].XxxHandler
func (s *App) Invoke(ctx context.Context, app, method string, in, rs interface{}) error {
	return s.invoke(ctx, app, method, cases.Title(language.English).String(app)+"Handler", in, rs)
}

// 更合理的，应该是通过port连接而不是定义app的map<?>
func (s *App) invoke(ctx context.Context, app, method, handler string, in, rs interface{}) error {
	var err error
	var conn *grpc.ClientConn

	//线上环境
	// if s.Conf.Environment == "--Environment--" {
	// 	conn, err = s.getConnFromRegistry(ctx, app)
	// } else {
	// 	conn, err = s.getConnFromPool(app)
	// }
	if err != nil {
		return err
	}
	return conn.Invoke(ctx, "/"+s.name+"."+app+"."+handler+"/"+method, in, rs)
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
