package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/starfork/go-slice"
)

type RequestData map[string]any

func (e RequestData) Set(key string, value any) RequestData {

	e[key] = value
	return e
}

type Request struct {
	h     http.Handler
	w     http.ResponseWriter
	req   *http.Request
	Data  RequestData
	Query url.Values
}

func New(h http.Handler, w http.ResponseWriter, req *http.Request) *Request {
	return &Request{
		h: h, w: w, req: req,
		Data:  make(RequestData),
		Query: req.URL.Query(),
	}
}

// 请求执行前打开请求
func (e *Request) Unpack() *Request {

	req := e.req
	var bd []byte
	if req.Method == "POST" || req.Method == "PUT" {
		bd, _ = io.ReadAll(req.Body)
		jsoniter.Unmarshal(bd, &e.Data)

	}
	//fmt.Println(req.Host)
	//ip := ClientIP(req)
	//e.Set("ip", ip)

	return e
}

// 解析完成之后会装req
func (e *Request) Pack() *Request {
	e.req.URL.RawQuery = e.Query.Encode()
	buf, _ := jsoniter.Marshal(e.Data)
	e.req.Body = io.NopCloser(bytes.NewBuffer(buf))

	return e
}

// 过滤不需要验证的url
func (e *Request) NoAuth(noAuth slice.Slice[string]) bool {
	return noAuth.Contains(e.req.URL.Path, func(key string) bool {
		return strings.HasPrefix(e.req.URL.Path, key)
	})
}

func (e *Request) Set(key, value string) {
	if e.Query == nil {
		e.Query = url.Values{}
	}
	if e.Data == nil {
		e.Data = make(RequestData)
	}
	e.Query.Set(key, value)
	e.Data.Set(key, value)
}

// []string{"X-Estate-Id"} 设置header部分变量
func (e *Request) SetHeader(keys []string) {
	for _, key := range keys {
		nkey := strings.ReplaceAll(strings.ToLower(key), "-", "_")
		v := e.req.Header.Get(key)
		if v != "" {
			e.Set(nkey, v)
		}
	}
}

func (e *Request) Next() {
	e.Pack()
	e.h.ServeHTTP(e.w, e.req)
}
