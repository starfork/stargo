package tagfile

//替换文件地址插件
//使用方法 1，在proto文件的字段中插入tag  `gorm:"tagfile:Nm;"`
//配置文件配置 fileserver

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/starfork/stargo/fileserver"
	"github.com/starfork/stargo/store"
	"github.com/starfork/stargo/util/ustring"
	"gorm.io/gorm"
)

type Plugin struct {
	fsc *fileserver.Config
}

func Register(db *gorm.DB, conf *store.Config) {

	p := &Plugin{
		fsc: conf.FileServer,
	}

	db.Callback().Query().After("gorm:find").Register("tagfile:after_query", p.AfterQuery)
	db.Callback().Update().Before("gorm:update").Register("tagfile:before_update", p.BeforeUpdate)

}

func (e *Plugin) Parse(str string) string {
	if str == "" {
		return ""
	}
	str = strings.ReplaceAll(str, "private://", e.fsc.PrivateUrl)
	str = strings.ReplaceAll(str, "public://", e.fsc.PublicUrl)
	return str
}
func (e *Plugin) Rebuild(str string) string {
	if str == "" {
		return ""
	}
	str = strings.ReplaceAll(str, e.fsc.PrivateUrl, "private://")
	str = strings.ReplaceAll(str, e.fsc.PublicUrl, "public://")
	return str
}

func (e *Plugin) AfterQuery(db *gorm.DB) {
	value := db.Statement.ReflectValue
	// if value.Kind() == reflect.Int64 {
	// 	return
	// }

	kind := value.Kind()
	if kind != reflect.Slice && kind != reflect.Struct {
		return
	}

	var tagFields []string
	//var moneyField []string

	for _, field := range db.Statement.Schema.Fields {
		//parseField, _ = field.TagSettings["TAGFILE"]---这样居然不行。golang这个变量作用域。。
		if f, ok := field.TagSettings["TAGFILE"]; ok {
			tagFields = append(tagFields, ustring.ToCamel(f))
		}
		// if f, ok := field.TagSettings["FMONEY"]; ok {
		// 	moneyField = append(moneyField, ustring.ToCamel(f))
		// }
	}
	if kind == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			//非指针类型的不能设置这些东东
			if reflect.Value(value.Index(i)).Kind() == reflect.Ptr {
				item := reflect.Value(value.Index(i)).Elem()
				e.SetPlugin(item, tagFields)
			}

		}
	} else if kind == reflect.Struct {
		item := reflect.Value(value)
		e.SetPlugin(item, tagFields)
	}

}

func (e *Plugin) SetPlugin(item reflect.Value, parseField []string) {

	if len(parseField) == 0 {
		return
	}
	//parseField = []string{"image_url", "image_url_list"}
	for _, field := range parseField {
		f := item.FieldByName(field)
		if (f != reflect.Value{}) { //字段没找到
			if f.CanSet() {
				f.SetString(e.Parse(f.String()))
			} else {
				fmt.Println("field can not set " + field)
			}
			//f.SetString(ReplacePrefix(f.String()))
		}
		// else {
		// 	if log != nil {
		// 		log.Info(context.TODO(), "field not found %s", field)
		// 	}
		// 	//log
		// }
	}
}

func (e *Plugin) BeforeUpdate(db *gorm.DB) {
	value := db.Statement.ReflectValue
	//fmt.Println(value.Kind())
	if value.Kind() != reflect.Struct {
		return
	}
	var parseField []string

	for _, field := range db.Statement.Schema.Fields {
		if f, ok := field.TagSettings["TAGFILE"]; ok {
			parseField = append(parseField, ustring.ToCamel(f))
		}
	}
	if len(parseField) == 0 {
		return
	}
	item := reflect.Value(value)
	for _, field := range parseField {
		f := item.FieldByName(field)
		if f.CanSet() {
			f.SetString(e.Rebuild(f.String()))
		} else {
			fmt.Println("field can not set " + field)
		}
	}
}
