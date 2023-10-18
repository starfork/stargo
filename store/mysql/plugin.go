package mysql

import (
	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/store/mysql/plugins/tagfile"
	"gorm.io/gorm"
)

var PluginsMap = map[string]func(db *gorm.DB, config *config.Config){
	"tagfile": tagfile.Register,
}

func RegisterPlugins(db *gorm.DB, conf *config.Config, plugins []string) {

	for _, name := range plugins {
		if f, ok := PluginsMap[name]; ok {
			f(db, conf)
		}
	}
}
