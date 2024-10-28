package api

import (
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

func WrapperSwagger(mux *http.ServeMux, fsys fs.FS) {
	mux.HandleFunc("/swagger/", swaggerFile)
	serveSwaggerUI(mux, fsys)
}

func serveSwaggerUI(mux *http.ServeMux, fsys fs.FS) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	fileServer := http.FileServer(http.FS(fsys))
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func swaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		log.Printf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	name := path.Join("openapiv2/proto/v1", p)
	log.Printf("Serving swagger-file: %s", name)
	http.ServeFile(w, r, name)

}
