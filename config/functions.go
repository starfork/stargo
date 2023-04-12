package config

import (
	"flag"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// LoadConfig config
func LoadConfig(config_file ...string) (*Config, error) {
	var configFile *string
	if len(config_file) > 0 {
		configFile = &config_file[0]
	} else {
		configFile = flag.String("c", "../config/debug.json", "config file path")
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

	decoder := jsoniter.NewDecoder(file)
	err = decoder.Decode(&conf)

	return conf, err
}
