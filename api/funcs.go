package api

import (
	"net"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Response(w http.ResponseWriter, code int, msg ...string) {
	type response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	res := response{
		Code:    code,
		Message: msg[0],
	}
	if len(msg) > 1 {
		res.Data = msg[1]
	}
	var ret, _ = jsoniter.Marshal(res)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(ret)

}

func Success(w http.ResponseWriter, msg ...string) {
	Response(w, http.StatusOK, msg...)
}
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return "unknow"
}
