package request

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

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
		log.Printf("oss-token: %s \n", token)
	}
	if tk == "" {
		return "", errors.New("token or authorization required")
	}
	return tk, nil
}
