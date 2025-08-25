package pm

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const EmptyString = ""

type Pm map[string]any

// 构造
func NewPm(data map[string]any) Pm {
	_pm := make(Pm)
	for k, v := range data {
		_pm[strings.ToLower(k)] = v
	}
	return _pm
}

// 设置参数
func (pm Pm) Set(key string, value any) Pm {
	if pm == nil {
		pm = make(Pm)
	}
	pm[strings.ToLower(key)] = value
	return pm
}

// 删除参数
func (pm Pm) Delete(keys ...string) Pm {
	if pm == nil {
		return pm
	}
	for _, key := range keys {
		delete(pm, strings.ToLower(key))
	}
	return pm
}

// 设置子 Pm
func (pm Pm) SetPm(key string, sub Pm) Pm {
	if pm == nil {
		pm = make(Pm)
	}
	pm[strings.ToLower(key)] = sub
	return pm
}

// 获取或创建子 Pm
func (pm Pm) SubPm(key string) Pm {
	if pm == nil {
		pm = make(Pm)
	}
	k := strings.ToLower(key)
	if v, ok := pm[k]; ok {
		if sub, ok := v.(Pm); ok {
			return sub
		}
	}
	sub := make(Pm)
	pm[k] = sub
	return sub
}

// 获取参数
func (pm Pm) Get(key string) any {
	if pm == nil {
		return nil
	}
	key = strings.ToLower(key)
	value, ok := pm[key]
	if !ok {
		return nil
	}
	return value
}

// --- 字符串获取 ---
func (pm Pm) GetString(key string) string {
	v := pm.Get(key)
	s, ok := v.(string)
	if ok {
		return s
	}
	return pm.toString(v)
}

// 严格字符串获取：只允许 string 类型
func (pm Pm) GetStringStrict(key string) string {
	v := pm.Get(key)
	if s, ok := v.(string); ok {
		return s
	}
	return EmptyString
}

// --- 数字获取 ---
func (pm Pm) GetInt(key string) int {
	if v, ok := pm.GetIntOk(key); ok {
		return v
	}
	return 0
}

func (pm Pm) GetIntOk(key string) (int, bool) {
	v := pm.Get(key)
	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		return int(x), true
	case string:
		i, err := strconv.Atoi(x)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

func (pm Pm) GetInt64(key string) int64 {
	if v, ok := pm.GetIntOk(key); ok {
		return int64(v)
	}
	return 0
}

func (pm Pm) GetUint32(key string) uint32 {
	if v, ok := pm.GetUint32Ok(key); ok {
		return v
	}
	return 0
}

func (pm Pm) GetUint32Ok(key string) (uint32, bool) {
	v := pm.Get(key)
	switch x := v.(type) {
	case uint32:
		return x, true
	case int:
		if x >= 0 {
			return uint32(x), true
		}
	case int64:
		if x >= 0 {
			return uint32(x), true
		}
	case float64:
		if x >= 0 {
			return uint32(x), true
		}
	case string:
		u, err := strconv.ParseUint(x, 10, 32)
		if err == nil {
			return uint32(u), true
		}
	}
	return 0, false
}

func (pm Pm) GetFloat64(key string) float64 {
	if v, ok := pm.GetFloat64Ok(key); ok {
		return v
	}
	return 0
}

func (pm Pm) GetFloat64Ok(key string) (float64, bool) {
	v := pm.Get(key)
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// --- Protobuf Any ---
// 这个使用体验并不好
// func (pm Pm) GetPbAny(key string) *structpb.Value {
// 	v := pm.Get(key)
// 	if value, err := structpb.NewValue(v); err == nil {
// 		return value
// 	}
// 	return nil
// }

// --- 内部工具 ---
func (pm Pm) toString(v any) string {
	if v == nil {
		return EmptyString
	}
	if s, ok := v.(string); ok {
		return s
	}
	bs, err := jsoniter.Marshal(v)
	if err != nil {
		return EmptyString
	}
	return string(bs)
}

// --- URL 编码 ---
func (pm Pm) EncodeURL() string {
	if pm == nil {
		return EmptyString
	}
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range pm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := pm.GetString(k)
		if v != EmptyString {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return EmptyString
	}
	return buf.String()[:buf.Len()-1]
}
