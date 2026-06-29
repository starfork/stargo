package cache

import (
	"github.com/starfork/stargo/cache"
	"github.com/starfork/stargo/store/mysql/plugins"
	"gorm.io/gorm"
)

type Plugin struct {
	c cache.Cache
}

func (p *Plugin) Name() string {
	return "cache"
}

func NewPlugin(c cache.Cache) *Plugin {
	return &Plugin{}
}

func Register(db *gorm.DB, conf plugins.Config) {
	p := &Plugin{}
	db.Callback().Query().After("gorm:find").Register("cache_after_query", p.After)
	db.Callback().Query().Before("gorm:find").Register("cache_before_query", p.Before)
}

func (e *Plugin) After(db *gorm.DB) {

}

func (e *Plugin) Before(db *gorm.DB) {

}
