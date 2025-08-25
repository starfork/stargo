package ustring

import (
	"fmt"
	"reflect"
	"strings"
)

func AnyToString(value any) string {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return "<nil>"
		}
		return AnyToString(v.Elem().Interface())
	case reflect.Map:
		var sb strings.Builder
		sb.WriteString("{")
		for i, key := range v.MapKeys() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s: %s", AnyToString(key.Interface()), AnyToString(v.MapIndex(key).Interface())))
		}
		sb.WriteString("}")
		return sb.String()
	case reflect.Slice:
		valueSlice := make([]any, v.Len())
		for i := range valueSlice {
			valueSlice[i] = AnyToString(v.Index(i).Interface())
		}
		return fmt.Sprintf("%v", valueSlice)
	default:
		return fmt.Sprintf("%v", value)
	}
}
