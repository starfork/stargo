package urange

import (
	"reflect"
	"strings"
)

// I2i 遍历对接rpc拿过来的数据。附加source的数据到target的对应字段里
// target 需要遍历的数据，source rpc过来的数据
// 目前仅支持string,uint类型
// (rs.Data, suser.Data, "Avatar,NickName", "", "UserId,Uid")
func I2i(target, source, targetKey, sourceKey, compareKey interface{}) {
	tKey := getTargetKey(targetKey)
	sKey := getSourceKey(tKey, sourceKey)
	cKey := getCompareKey(compareKey)
	sTarget := reflect.ValueOf(target)
	sSource := reflect.ValueOf(source)

	for i := 0; i < sTarget.Len(); i++ {
		so1 := sTarget.Index(i).Elem()
		f1 := so1.FieldByName(cKey[0]).Interface() //compare key 1
		for j := 0; j < sSource.Len(); j++ {
			de1 := sSource.Index(j).Elem()
			f2 := sSource.Index(j).Elem().FieldByName(cKey[1]).Interface() //compare key 2
			if f1 == f2 {
				for k := 0; k < len(tKey); k++ {
					f3 := de1.FieldByName(sKey[k])

					sf := so1.FieldByName(tKey[k])
					dataKind := getKind(f3)
					switch dataKind {
					case reflect.String:
						sf.SetString(f3.String())
					case reflect.Uint:
						//fmt.Println(f3.Uint())
						sf.SetUint(uint64(f3.Uint()))
					}
				}
			}
		}
		//rs = append(rs, so.Index(i))
	}
}

// I2i 遍历对接rpc拿过来的数据。附加source的数据到target的对应字段里
// target 需要遍历的数据，source rpc过来的数据
// 目前仅支持string,uint类型
// targetKey 目标key，sourceKey 资源key。如果sourceKey为空，则使用targetKey。
// compareKey对比key
// subkey二级struct名字
// urange.I2iKey(data, suser.Data, "Avatar,NickName", "", "FriendId,Uid", "worker")
func I2iKey(target, source, targetKey, sourceKey, compareKey interface{}, subkey string) {
	tKey := getTargetKey(targetKey)
	sKey := getSourceKey(tKey, sourceKey)
	cKey := getCompareKey(compareKey)
	sTarget := reflect.ValueOf(target)
	sSource := reflect.ValueOf(source)

	for i := 0; i < sTarget.Len(); i++ {
		iTarget := sTarget.Index(i).Elem()
		tField := iTarget.FieldByName(cKey[0]).Interface() //compare key 1
		for j := 0; j < sSource.Len(); j++ {
			de1 := sSource.Index(j).Elem()
			//fmt.Println(cKey)
			sField := de1.FieldByName(cKey[1]).Interface() //compare key 2

			if tField == sField {
				sf := reflect.Indirect(iTarget).FieldByName(subkey)
				if sf.Kind() != reflect.Ptr {
					//log warn
					break
				}
				tmp := reflect.New(sf.Type().Elem())
				t := reflect.Indirect(tmp)
				for k := 0; k < len(tKey); k++ {
					f3 := de1.FieldByName(sKey[k])
					sf := t.FieldByName(tKey[k])
					dataKind := getKind(f3)
					switch dataKind {
					case reflect.String:
						sf.SetString(f3.String())
					case reflect.Uint:
						sf.SetUint(uint64(f3.Uint()))
					}
				}
				sf.Set(tmp)
				iTarget.FieldByName(subkey).Set(sf)
			}
		}
		//rs = append(rs, so.Index(i))
	}
}

// getTargetKey 获取target 的key。估计后面都不会用string的slice了。
func getTargetKey(key interface{}) []string {
	if op, ok := key.(string); ok {
		return strings.Split(op, ",")
	}
	return key.([]string)
}

// getSourceKey 获取Source的key
func getSourceKey(targetKey []string, key interface{}) []string {
	if op, ok := key.(string); ok {
		if op == "" {
			return targetKey
		}
		keys := strings.Split(op, ",")
		if len(keys) != len(targetKey) {
			return targetKey
		}
		return keys
	}
	return key.([]string)
}

// getCompareKey 解析对比Key。
// 如果没逗号，表示两个key相同
// 如果有逗号。第一个表示target的key，第二个是source key
func getCompareKey(key interface{}) []string {
	if op, ok := key.(string); ok {
		tmp := strings.Split(op, ",")
		if len(tmp) == 1 {
			return []string{op, op}
		}
		return tmp
	}
	return key.([]string)
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
