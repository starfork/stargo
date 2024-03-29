package mysql

import (
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/store/mysql/plugins/tagfile"
	"github.com/starfork/stargo/store/mysql/plugins/unmarshaler"
	"gorm.io/gorm"
)

var PluginsMap = map[string]func(db *gorm.DB, config *store.Config){
	"tagfile":     tagfile.Register,
	"unmarshaler": unmarshaler.Register,
}

func RegisterPlugins(db *gorm.DB, conf *store.Config, plugins []string) {

	for _, name := range plugins {
		if f, ok := PluginsMap[name]; ok {
			f(db, conf)
		}
	}
}
