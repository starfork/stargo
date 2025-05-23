package api

import (
	"io/fs"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
)

type Config struct {
	App             string
	Port            string
	Registry        *naming.Config
	DiaOpts         []grpc.DialOption
	SwgFs           fs.FS
	Wrapper         func(http.Handler) http.Handler
	SwaggerRoute    string
	SwaggerUIPrefix string
	SMOpts          []runtime.ServeMuxOption

	MuxHandler map[string]func(w http.ResponseWriter, r *http.Request)
}
