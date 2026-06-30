package api

import (
	"net/http"
	"strings"
)

// CORSConfig holds CORS middleware configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

var DefaultCORSConfig = CORSConfig{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With", "X-Request-Id"},
	ExposedHeaders:   []string{"X-Request-Id"},
	AllowCredentials: false,
	MaxAge:           86400,
}

// CORSWrapper wraps an http.Handler with CORS headers
func CORSWrapper(cfg CORSConfig) func(http.Handler) http.Handler {
	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = DefaultCORSConfig.AllowedMethods
	}
	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = DefaultCORSConfig.AllowedHeaders
	}
	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = DefaultCORSConfig.AllowedOrigins
	}

	allowedMethods := strings.Join(cfg.AllowedMethods, ", ")
	allowedHeaders := strings.Join(cfg.AllowedHeaders, ", ")
	exposedHeaders := strings.Join(cfg.ExposedHeaders, ", ")
	allowCredentials := "false"
	if cfg.AllowCredentials {
		allowCredentials = "true"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if cfg.isOriginAllowed(origin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", allowCredentials)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (c *CORSConfig) isOriginAllowed(origin string) bool {
	if len(c.AllowedOrigins) == 1 && c.AllowedOrigins[0] == "*" {
		return true
	}
	for _, o := range c.AllowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}
