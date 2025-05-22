package plugins

import "gorm.io/gorm"

type Config map[string]any

type Plugin interface {
	gorm.Plugin
}
