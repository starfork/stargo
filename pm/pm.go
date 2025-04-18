package pm

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	NULL = ""
)

type Pm map[string]any

func NewPm(data map[string]any) Pm {
	_pm := make(Pm)
	for k, v := range data {
		_pm[strings.ToLower(k)] = v
	}
	return _pm
}

// 设置参数
func (pm Pm) Set(key string, value any) Pm {
	key = strings.ToLower(key)
	pm[key] = value
	return pm
}

func (pm Pm) Delete(keys ...string) Pm {
	for _, key := range keys {
		delete(pm, key)
	}
	return pm
}
func (pm Pm) SetPm(key string, f func(b Pm)) Pm {
	key = strings.ToLower(key)
	_pm := make(Pm)
	f(_pm)
	pm[key] = _pm
	return pm
}

// 获取参数，同 GetString()
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

// 获取参数转换string
func (pm Pm) GetString(key string) string {
	value := pm.Get(key)
	v, ok := value.(string)
	if !ok {
		return pm.toString(value)
	}
	return v
}

func (pm Pm) GetPbAny(key string) *structpb.Value {
	v := pm.Get(key)
	if value, err := structpb.NewValue(v); err == nil {
		return value
	}
	return nil
}

// 获取float64。其他float再转一下
func (pm Pm) GetFloat64(key string) float64 {
	value := pm.GetString(key)

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return v
}
func (pm Pm) GetUint32(key string) uint32 {
	return uint32(pm.GetInt(key))
}

// 获取int。其他int自己再转一下
func (pm Pm) GetInt(key string) int {
	value := pm.GetString(key)

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v
}
func (pm Pm) GetInt64(key string) int64 {
	return int64(pm.GetInt(key))
}

func (pm Pm) toString(v any) (str string) {

	if v == nil {
		return NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = jsoniter.Marshal(v); err != nil {
		return NULL
	}
	return string(bs)
}

// 编码URL
func (pm Pm) EncodeURL() string {
	if pm == nil {
		return NULL
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
		if v := pm.GetString(k); v != NULL {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return NULL
	}
	return buf.String()[:buf.Len()-1]

}
