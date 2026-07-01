package config

import (
	"fmt"
	"time"

	"github.com/starfork/stargo/api"
	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/secrets"
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
	Timezone   string
	Timeformat string

	Server *server.Config
	Api    *api.Config

	Store map[string]*store.Config

	Log      *logger.Config
	Broker   *broker.Config
	Registry *naming.Config
	Tracer   *tracer.Config
	Secret   *secrets.Config

	//Jwt *JwtConfig
}

var DefaultConfig = &Config{
	Server: server.DefaultConfig,
	Store:  make(map[string]*store.Config),
	Log:    &logger.Config{},
	Env:    ENV_DEV,
}

func (c *Config) Validate() error {
	if c.Env == "" {
		return fmt.Errorf("config: env is required")
	}
	if c.Server == nil {
		return fmt.Errorf("config: server config is required")
	}
	if c.Server.Addr == "" {
		return fmt.Errorf("config: server.addr is required")
	}
	if c.Server.ShutdownTimeout <= 0 {
		c.Server.ShutdownTimeout = 30 * time.Second
	}
	for name, sc := range c.Store {
		if sc == nil {
			continue
		}
		if sc.Host == "" && sc.DSN == "" {
			return fmt.Errorf("config: store[%s]: host or dsn is required", name)
		}
	}
	return nil
}

// type JwtConfig struct {
// 	PublicKey  string
// 	PrivateKey string
// }
