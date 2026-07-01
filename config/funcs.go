package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/starfork/stargo/secrets"
	"gopkg.in/yaml.v3"
)

// 一般的，都会实现自己的config，这里当作参考
// LoadConfig config
func LoadConfig(config_file ...string) (*Config, error) {
	var configFile *string
	if len(config_file) > 0 {
		configFile = &config_file[0]
	} else {
		configFile = flag.String("c", "../config/debug.yaml", "config file path")
	}
	flag.Parse()
	return ParseConfig(*configFile)
}

// ParseConfig 解析
func ParseConfig(f string) (*Config, error) {

	conf := &Config{}
	file, err := os.Open(f)
	if err != nil {
		return conf, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&conf)
	if err != nil {
		return conf, err
	}

	if conf.Secret != nil && conf.Secret.Driver != "" {
		sm, err := secrets.New(conf.Secret.Driver, conf.Secret)
		if err != nil {
			return conf, fmt.Errorf("secret manager init: %w", err)
		}
		defer sm.Close()

		if err := secrets.ResolveAll(sm, conf); err != nil {
			return conf, fmt.Errorf("secret resolve: %w", err)
		}
	}

	if err := conf.Validate(); err != nil {
		return conf, fmt.Errorf("config validation failed: %w", err)
	}
	return conf, nil
}
