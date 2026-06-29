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

func RegisterPlugins(db *gorm.DB, pluginCfgs map[string]map[string]any) {

	//db.Use()

	for name, cfg := range pluginCfgs {
		if f, ok := PluginsMap[name]; ok {
			f(db, plugins.Config(cfg))
		}
	}
}
