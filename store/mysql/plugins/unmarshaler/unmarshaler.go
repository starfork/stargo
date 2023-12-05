package unmarshaler

/**
自动调用解析
在repository中，统一规范字段返解操作
在repository文件中，定义一个Unmarshal函数，在里面写具体里逻辑即可

type Addon struct {
}
func (p *Addon) Unmarshal() {
	//for test
	p.Addon.Author = p.Addon.Author + "___decode"
}
*/

import (
	"reflect"

	"github.com/starfork/stargo/config"
	"gorm.io/gorm"
)

type Plugin struct {
}

func Register(db *gorm.DB, conf *config.Config) {
	p := &Plugin{}
	db.Callback().Query().After("gorm:find").Register("tagfile:after_query", p.AfterQuery)
}

// 接口定义
type Unmarshaler interface {
	Unmarshal()
}

func (e *Plugin) AfterQuery(db *gorm.DB) {
	value := db.Statement.ReflectValue
	kind := value.Kind()
	if kind != reflect.Slice && kind != reflect.Struct {
		return
	}

	if kind == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			//非指针类型的不能设置这些东东
			if reflect.Value(value.Index(i)).Kind() == reflect.Ptr {
				item := reflect.Value(value.Index(i)).Elem()
				e.unmarshal(item)
			}

		}
	} else if kind == reflect.Struct {
		item := reflect.Value(value)
		e.unmarshal(item)
	}

}

func (e *Plugin) unmarshal(item reflect.Value) {
	it := reflect.TypeOf((*interface{ Unmarshal() })(nil)).Elem()
	addr := item.Addr()

	if addr.Type().Implements(it) {
		addr.MethodByName("Unmarshal").Call(nil)
	}
}
