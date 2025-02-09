package api

import (
	"net/http"
	"strings"
)

func Cros(w http.ResponseWriter, req *http.Request, headers ...map[string]string) {
	req.Header.Set("Grpc-Metadata-G-Method", req.Method)
	req.Header.Set("Grpc-Metadata-IP", ClientIP(req))
	req.Header.Set("Grpc-Metadata-Host", req.Host)
	req.Header.Set("Grpc-Metadata-Token", req.Header.Get("Access-Token"))

	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set("Grpc-Metadata-"+k, v)
		}
	}

	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			headers := []string{"Content-Type", "Access-Token", "Origin", "Authorization", "Accept", "Content-Type", "X-Requested-With"}
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			return
		}
	}
}
