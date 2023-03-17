package request

import (
	"errors"
	"net/http"
	"strings"

	"public/pkg/util/crypt/jwt"

	"github.com/golang/glog"
)

// Token 检查header获取token
func Token(req *http.Request) (*jwt.Options, error) {
	//var tk string
	var rs *jwt.Options
	var err error

	tk, err := GetToken(req)
	if err != nil {
		return nil, errors.New("token or authorization required")
	}
	//glog.Infof("access-token: %s ", tk)

	if rs, err = jwt.Parse(tk); err != nil {
		return nil, errors.New("token invalid:" + err.Error())
	}
	return rs, nil
}

func GetToken(req *http.Request) (string, error) {
	var tk string
	// fmt.Println(req.URL.Query())
	if token, ok := req.Header["Access-Token"]; ok && len(token) > 0 {
		tk = token[0]
	} else if token, ok := req.Header["Authorization"]; ok && len(token) > 0 {
		tk = strings.Split(token[0], " ")[1]
	} else if token, ok := req.URL.Query()["access-token"]; ok && len(token) > 0 {
		tk = token[0]
	} else if token, ok := req.URL.Query()["oss-token"]; ok && len(token) > 0 {
		// todo oss临时访问token。
		// access-token过长。 go的 http.get 会出问题
		glog.Infof("oss-token: %s ", token)
	}
	if tk == "" {
		return "", errors.New("token or authorization required")
	}
	return tk, nil
}
