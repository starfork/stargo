package config

import (
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/store"
)

var (
	ENV_DEBUG      = "debug"
	ENV_PRODUCTION = "production"
)

// 公共配置模板
type Config struct {
	Env string
	Org string
	//ServerName string //服务名称--4-11改。通过app启动设置
	RpcServer  *RpcServer
	Rpc        map[string]*RpcServer
	HttpServer *HttpServer

	Timezome   string //时区设置
	Timeformat string

	Store map[string]*store.Config

	Log      *logger.Config
	Broker   *broker.Config
	Registry *naming.Config
}

type RpcServer struct {
	Entry string
	Name  string
	Host  string
	Port  string
	Auth  string //[keyfilepath]:[key]:
}

type HttpServer struct {
	Host string
	Port string
}
