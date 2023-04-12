package plugin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/util/ustring"
	"gorm.io/gorm"
)

type Plugin struct {
	//file server config
	fsc *config.FileServerConfig
	//log *config.Config
}

func New(conf *config.ServerConfig) *Plugin {
	return &Plugin{
		fsc: conf.FileServer,
	}
}

// 接口定义
type Unmarshaler interface {
	Unmarshal()
}

func (e *Plugin) ReplacePrefix(str string) string {
	if str == "" {
		return ""
	}

	str = strings.ReplaceAll(str, "private://", e.fsc.PrivateStaticUrl)
	str = strings.ReplaceAll(str, "public://", e.fsc.PublicStaticUrl)
	return str
}
func (e *Plugin) RebuildPrefix(str string) string {
	if str == "" {
		return ""
	}
	str = strings.ReplaceAll(str, e.fsc.PrivateStaticUrl, "private://")
	str = strings.ReplaceAll(str, e.fsc.PublicStaticUrl, "public://")
	return str
}

func (e *Plugin) AfterQuery(db *gorm.DB) {

	value := db.Statement.ReflectValue
	if value.Kind() == reflect.Int64 {
		return
	}

	var imgParseField []string
	//var moneyField []string

	for _, field := range db.Statement.Schema.Fields {
		//parseField, _ = field.TagSettings["IMGPARSE"]---这样居然不行。golang这个变量作用域。。
		if f, ok := field.TagSettings["IMGPARSE"]; ok {
			imgParseField = append(imgParseField, ustring.ToCamel(f))
		}
		// if f, ok := field.TagSettings["FMONEY"]; ok {
		// 	moneyField = append(moneyField, ustring.ToCamel(f))
		// }
	}
	if value.Kind() == reflect.Slice {

		for i := 0; i < value.Len(); i++ {
			//非指针类型的不能设置这些东东
			if reflect.Value(value.Index(i)).Kind() == reflect.Ptr {
				item := reflect.Value(value.Index(i)).Elem()
				unmarshal(item)
				e.SetOssImg(item, imgParseField)
			}

		}
	} else if value.Kind() == reflect.Struct {
		item := reflect.Value(value)
		unmarshal(item)
		e.SetOssImg(item, imgParseField)
	}

}

func unmarshal(item reflect.Value) {
	it := reflect.TypeOf((*interface{ Unmarshal() })(nil)).Elem()
	addr := item.Addr()
	if addr.Type().Implements(it) {
		addr.MethodByName("Unmarshal").Call(nil)

	}
}

func (e *Plugin) SetOssImg(item reflect.Value, parseField []string) {

	if len(parseField) == 0 {
		return
	}
	//parseField = []string{"image_url", "image_url_list"}
	for _, field := range parseField {
		f := item.FieldByName(field)
		if (f != reflect.Value{}) { //字段没找到
			if f.CanSet() {
				f.SetString(e.ReplacePrefix(f.String()))
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
		if f, ok := field.TagSettings["IMGPARSE"]; ok {
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
			f.SetString(e.RebuildPrefix(f.String()))
		} else {
			fmt.Println("field can not set " + field)
		}
	}
}
