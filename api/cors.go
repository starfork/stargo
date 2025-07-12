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
		"Access-Token", "Access-Fp", "Accept-Language", "Access-Device",
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
		META_DEVICE: req.Header.Get("Access-Device"),
	}
	for k, v := range inHeaders {
		req.Header.Set("Grpc-Metadata-"+k, v)
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set("Grpc-Metadata-"+k, v)
		}
	}

	/**
	 * 更多时候推荐直接nginx去处理
	 * 如果这里和nginx有冲突。可以按照如下方案处理
		proxy_hide_header Access-Control-Allow-Origin;
		proxy_hide_header Access-Control-Allow-Methods;
		proxy_hide_header Access-Control-Allow-Headers;

		add_header 'Access-Control-Allow-Origin' '*' always;
		add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
		add_header 'Access-Control-Allow-Headers' 'Authorization,Content-Type,Accept,Origin,X-Requested-With,Access-Token,Access-Fp' always;
		#这里需要cdn如果有启用http2的话会有类似cloudflare 520的问题。需要关闭掉
		#cloudflare的处理方式-》“速度” > “优化” > “协议优化”中禁用与源站的 HTTP/2 连接。
		if ($request_method = 'OPTIONS') {
			add_header 'Access-Control-Allow-Origin' '*' always;
			add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
			add_header 'Access-Control-Allow-Headers' 'Authorization,Content-Type,Accept,Origin,X-Requested-With,Access-Token,Access-Fp' always;
			add_header 'Content-Length' 0;
			add_header 'Content-Type' 'text/plain; charset=UTF-8';
			return 200;
		}
	*/
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowHeaders, ","))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(AllowMethods, ","))
			return
		}
	}
}
