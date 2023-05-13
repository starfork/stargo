package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"strings"

	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/starfork/stargo/util/crypt/jwt"
	"github.com/starfork/stargo/util/gslice"
	"github.com/starfork/stargo/util/request"
)

type RequestData map[string]interface{}

func (e RequestData) Set(key string, value interface{}) RequestData {

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

// token解析
func (e *Request) ParseToken() (*jwt.Options, error) {
	tk, err := request.GetToken(e.req)
	if err != nil {
		return nil, err
	}
	return jwt.Parse(tk)

}

// 请求执行前打开请求
func (e *Request) Unpack() {

	req := e.req
	var bd []byte
	if req.Method == "POST" || req.Method == "PUT" {
		bd, _ = io.ReadAll(req.Body)
		jsoniter.Unmarshal(bd, &e.Data)

	}
	glog.Infof("%s %s %s %s ", req.Method, req.URL, string(bd), req.Header.Get("Access-Token"))

}

// 解析完成之后会装req
func (e *Request) Pack() {
	e.req.URL.RawQuery = e.Query.Encode()
	buf, _ := jsoniter.Marshal(e.Data)
	e.req.Body = io.NopCloser(bytes.NewBuffer(buf))

}

// 过滤不需要验证的url
func (e *Request) NoAuth(noAuth gslice.Slice[string]) bool {
	return noAuth.Contains(e.req.URL.Path, func(key string) bool {
		return strings.Contains(e.req.URL.Path, key)
	})

}

func (e *Request) Set(key, value string) {

	e.Query.Set(key, value)
	e.Data.Set(key, value)
}

func (e *Request) Next() {
	e.Pack()
	e.h.ServeHTTP(e.w, e.req)
}
