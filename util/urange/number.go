package urange

import (
	"fmt"
	"reflect"
)

// Ruint32 遍历给定对象返回一个uint32 slice
//
//		type Source struct { Id uint32 }
//		var rs []*Source
//		rs = append(rs, &Source{Id: 111})
//		rs = append(rs, &Source{Id: 222})
//	 uid := Ruint32(rs, "Id")
//
// 对象的key区分大小写。返回结果已经去重复
func Ruint32(source any, key string) []uint32 {
	var uid []uint32
	for _, v := range rdi(source, key) {
		uid = append(uid, v.(uint32))
	}
	return uid
}

// 拼接成 'bbb','xxx'
func Rstring(source any, key string) string {
	var uid string
	for _, v := range rdi(source, key) {
		uid += fmt.Sprintf("%s,", v.(string))
	}
	return uid
}

// 拼接成 ['xxx','xxx']
func Rstrings(source any, key string) []string {
	var uid []string
	for _, v := range rdi(source, key) {
		uid = append(uid, v.(string))
	}
	return uid
}

// 去重
func Ruint64(source any, key string) []uint64 {
	var uid []uint64
	for _, v := range rdi(source, key) {
		uid = append(uid, v.(uint64))
	}
	return uid
}

// ROuint32 原始的，没有去重复的遍历数据
func ROuint32(source any, key string) []uint32 {
	var uid []uint32
	for _, v := range ri(source, key) {
		uid = append(uid, v.(uint32))
	}
	return uid
}

// 原始
func ROuint64(source any, key string) []uint64 {
	var uid []uint64
	for _, v := range ri(source, key) {
		uid = append(uid, v.(uint64))
	}
	return uid
}

// ri 原始interface
func ri(source any, key string) []any {
	var uid []any
	val := reflect.ValueOf(source)
	for i := 0; i < val.Len(); i++ {
		k := val.Index(i).Elem()
		f := k.FieldByName(key).Interface()
		uid = append(uid, f)
	}
	return uid
}

// ri range duplicate interface 原始去重复interface
func rdi(source any, key string) []any {
	//var uid []any
	val := reflect.ValueOf(source)
	uid := make([]any, 0, val.Len())
	temp := map[any]struct{}{}
	for i := 0; i < val.Len(); i++ {
		f := val.Index(i).Elem().FieldByName(key).Interface()
		if _, ok := temp[f]; !ok {
			temp[f] = struct{}{}
			uid = append(uid, f)
		}
	}
	return uid
}
