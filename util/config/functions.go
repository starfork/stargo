package config

import (
	"strconv"
	"strings"
)

// Get config value
func (e *Config) Get(name string, deft ...string) string {
	if v, ok := e.val[name]; ok {
		return v
	}
	if len(deft) > 0 {
		return deft[0]
	}
	return ""
}

func (e *Config) GetAll() Val {
	return e.val
}

func (e *Config) GetInt(name string, deft ...int) int {
	if v, ok := e.val[name]; ok {
		i, _ := strconv.Atoi(v)
		return i
	}
	if len(deft) > 0 {
		return deft[0]
	}
	return 0
}
func (e *Config) GetUint32(name string, deft ...int) uint32 {
	return uint32(e.GetInt(name, deft...))
}
func (e *Config) GetUint64(name string, deft ...int) uint64 {
	return uint64(e.GetInt(name, deft...))
}
func (e *Config) GetInt64(name string, deft ...int) int64 {
	return int64(e.GetInt(name, deft...))
}

// 多行的，用竖线隔开的配置。 "abc|def"=>map[abc:def]
func (e *Config) getSlice(name string, sep string) map[string]string {
	mp := make(map[string]string)
	if v, ok := e.val[name]; ok {
		for _, vv := range strings.Split(v, "\n") {
			vvv := strings.Split(vv, sep)
			if len(vvv) > 1 {
				mp[vvv[0]] = vvv[1]
			}
		}
	}
	return mp
}

//sortMapKey golang的map range是无序的，所以需要自己排序
// func sortMapKey(mp map[string]string) []string {
// 	var newMap = make([]string, 0)
// 	for k := range mp {
// 		newMap = append(newMap, k)
// 	}
// 	sort.Strings(newMap)
// 	return newMap

// }

func (e *Config) GetStepInt(value int, name string) int {
	i, _ := strconv.Atoi(e.GetStep(value, name))
	return i
}

// 对于多行策略，如下
// 100|200,xxx,yyy
// 300|456,sdg,234
// 900|87s,max.sdfs.sdg
// 传递，301-900，返回“456,sdg,234”，传递，大于901返回“|87s,max.sdfs.sdg”
// 传递 0-99，返回空
func (e *Config) GetStep(value int, name string) string {
	mp := e.getSlice(name, "|")
	//fmt.Println(mp)
	tmp := ""
	for k, v := range mp {
		cKey, _ := strconv.Atoi(k)
		if value < cKey {
			return tmp
		}
		tmp = v
	}
	return tmp
}

//===检查某些元素是否在配置文件中==========

func (e *Config) Exist(key, value string, sep ...string) bool {

	s := ","
	if len(sep) > 0 {
		s = sep[0]
	}

	rs := strings.Split(e.Get(key), s)
	for _, v := range rs {
		if v == value {
			return true
		}
	}

	return false
}

// ExistUint32
func (e *Config) ExistUint32(key string, value uint32, sep ...string) bool {
	v := strconv.Itoa(int(value))
	return e.Exist(key, v, sep...)
}
