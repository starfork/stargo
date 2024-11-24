package api

import (
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

func (e *Api) WrapperSwagger(mux *http.ServeMux) {
	route := "/swagger/"
	if e.conf.SwaggerRoute != "" {
		route = e.conf.SwaggerRoute
	}
	mux.HandleFunc(route, e.swaggerFile(route))
	e.serveSwaggerUI(mux, e.conf.SwgFs)
}

func (e *Api) serveSwaggerUI(mux *http.ServeMux, fsys fs.FS) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	fileServer := http.FileServer(http.FS(fsys))
	prefix := "/swagger-ui/"
	if e.conf.SwaggerUIPrefix != "" {
		prefix = e.conf.SwaggerUIPrefix
	}
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func (e *Api) swaggerFile(prefix string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			log.Printf("Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, prefix)
		name := path.Join("openapiv2/proto/v1", p)
		log.Printf("Serving swagger-file: %s", name)
		http.ServeFile(w, r, name)
	}
}
