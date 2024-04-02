package config

import (
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/fileserver"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
)

var (
	ENV_DEBUG      = "debug"
	ENV_PRODUCTION = "production"
)

// type Config struct {
// 	Deploy  string
// 	Monitor bool
// 	Base    *ServerConfig //如果各个服务么有单独设置，则公用
// 	Server  map[string]*ServerConfig
// }

// 公共配置模板
type Config struct {
	Environment string
	Org         string
	//ServerName string //服务名称--4-11改。通过app启动设置
	Port string //服务端口
	Xds  bool   //是否是xds类型

	Timezome   string //时区设置
	Timeformat string

	Mysql  *store.Config
	Redis  *store.Config
	Mongo  *store.Config
	Sqlite *store.Config

	FileServer *fileserver.Config
	Log        *logger.Config
	Broker     *broker.Config
	Registry   *naming.Config

	//RpcServer map[string]*Server //rpc server
}

// log

//Rpc Server

type Server struct {
	Entry string
	Name  string
	Host  string
	Port  string
	Auth  string //[keyfilepath]:[key]:
}
