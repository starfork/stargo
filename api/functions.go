package api

import (
	"net/http"

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
