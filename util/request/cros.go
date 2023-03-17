package request

import (
	"net/http"
	"strings"
)

func Cros(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Grpc-Metadata-G-Method", req.Method)
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			headers := []string{"Content-Type", "Accept"}
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			return
		}
	}
}
