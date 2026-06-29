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

	"github.com/starfork/stargo/store/mysql/plugins"
	"gorm.io/gorm"
)

type Plugin struct {
}

func Register(db *gorm.DB, conf plugins.Config) {
	p := &Plugin{}
	db.Callback().Query().After("gorm:find").Register("unmarshaler_after_query", p.After)
	//db.Callback().Create().Before("gorm:create").Register("unmarshaler_after_create", p.Before)
	//db.Callback().Update().Before("gorm:update").Register("unmarshaler_after_update", p.Before)
}

func (e *Plugin) After(db *gorm.DB) {

	value := db.Statement.ReflectValue
	kind := value.Kind()
	if kind != reflect.Slice && kind != reflect.Struct {
		return
	}

	switch kind {
	case reflect.Slice:
		for i := range value.Len() {
			//非指针类型的不能设置这些东东
			if reflect.Value(value.Index(i)).Kind() == reflect.Ptr {
				item := reflect.Value(value.Index(i)).Elem()
				e.unmarshal(db, item)
			}

		}
	case reflect.Struct:
		item := reflect.Value(value)
		e.unmarshal(db, item)
	}

}

func (e *Plugin) Before(db *gorm.DB) {
	value := db.Statement.ReflectValue
	kind := value.Kind()
	if kind != reflect.Slice && kind != reflect.Struct {
		return
	}

	switch kind {
	case reflect.Slice:
		for i := range value.Len() {
			//非指针类型的不能设置这些东东
			if reflect.Value(value.Index(i)).Kind() == reflect.Ptr {
				item := reflect.Value(value.Index(i)).Elem()
				e.marshal(item)
			}

		}
	case reflect.Struct:
		item := reflect.Value(value)
		e.marshal(item)
	}

}

// 可支持参数
func (e *Plugin) unmarshal(db *gorm.DB, item reflect.Value) {
	addr := item.Addr()
	m := addr.MethodByName("Unmarshal")
	if !m.IsValid() {
		return
	}

	mt := m.Type()
	n := mt.NumIn()

	val, _ := db.Get("Sence")
	if val != nil {
		m.Call([]reflect.Value{reflect.ValueOf(val)})
		return
	}

	if n == 0 {
		m.Call(nil)
		return
	}
	args := make([]reflect.Value, n)
	for i := range n {
		t := mt.In(i)
		args[i] = reflect.Zero(t)
	}

	m.Call(args)
}

func (e *Plugin) marshal(item reflect.Value) {
	it := reflect.TypeOf((*interface{ Marshal() })(nil)).Elem()
	addr := item.Addr()

	if addr.Type().Implements(it) {
		addr.MethodByName("Marshal").Call(nil)
	}
}
