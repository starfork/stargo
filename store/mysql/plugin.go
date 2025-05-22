package mysql

import (
	"github.com/starfork/stargo/store/mysql/plugins"
	"github.com/starfork/stargo/store/mysql/plugins/tagfile"
	"github.com/starfork/stargo/store/mysql/plugins/unmarshaler"
	"gorm.io/gorm"
)

var PluginsMap = map[string]func(db *gorm.DB, config plugins.Config){
	"tagfile":     tagfile.Register,
	"unmarshaler": unmarshaler.Register,
}

func RegisterPlugins(db *gorm.DB, plugins map[string]plugins.Config) {

	//db.Use()

	for name, config := range plugins {
		if f, ok := PluginsMap[name]; ok {
			f(db, config)
		}
	}
}
