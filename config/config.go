package config

import (
	"github.com/starfork/stargo/api"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/filemanager"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/server"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/tracer"
)

var (
	ENV_DEV        = "dev"        //本地测试环境
	ENV_DOCKER     = "docker"     //docker模式
	ENV_PRODUCTION = "production" //正式环境
)

// 公共配置模板
type Config struct {
	Env        string
	Timezome   string //时区设置
	Timeformat string
	//Cm         []string // config manager。目前只处理基于etcd，所以这里是etcd的地址

	Server *server.Config
	Api    *api.Config

	Store map[string]*store.Config

	Log         *logger.Config
	Broker      *broker.Config
	Registry    *naming.Config
	Tracer      *tracer.Config
	Filemanager *filemanager.Config
}

var DefaultConfig = &Config{
	Server:   server.DefaultConfig,
	Store:    make(map[string]*store.Config),
	Log:      &logger.Config{},
	Broker:   &broker.Config{},
	Registry: &naming.Config{},
	//Tracer:   &tracer.Config{},
}
