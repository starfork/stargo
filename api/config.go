package api

import (
	"io/fs"
	"net/http"

	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
)

type Config struct {
	App      string
	Port     string
	Registry *naming.Config
	DiaOpts  []grpc.DialOption
	SwgFs    fs.FS
	Wrapper  func(http.Handler) http.Handler
}
