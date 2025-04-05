package api

import (
	"net/http"
	"strings"
)

const (
	META_FP     = "Stargo-Fp"
	META_TOKEN  = "Stargo-Token"
	META_METHOD = "Stargo-Method"
	META_IP     = "Stargo-IP"
	META_HOST   = "Stargo-Host"
	META_LANG   = "Stargo-Lang"
	META_DEVICE = "Stargo-Device"
)

var (
	AllowHeaders = []string{"Content-Type", "Origin", "Authorization", "Content-Type", "X-Requested-With",
		"Accept", "Access-Control-Allow-Credentials",
		"Access-Token", "Access-Fp", "Accept-Language", "Accept-Device",
	}
	AllowMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
)

func Cros(w http.ResponseWriter, req *http.Request, headers ...map[string]string) {
	inHeaders := map[string]string{
		META_METHOD: req.Method,
		META_IP:     ClientIP(req),
		META_HOST:   req.Host,
		META_TOKEN:  req.Header.Get("Access-Token"),
		META_FP:     req.Header.Get("Access-Fp"),
		META_LANG:   req.Header.Get("Accept-Language"),
		META_DEVICE: req.Header.Get("Accept-Device"),
	}
	for k, v := range inHeaders {
		req.Header.Set("Grpc-Metadata-"+k, v)
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set("Grpc-Metadata-"+k, v)
		}
	}

	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowHeaders, ","))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(AllowMethods, ","))
			return
		}
	}
}
