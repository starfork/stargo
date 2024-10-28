package api

import (
	"io/fs"

	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
)

type Config struct {
	App      string
	Port     string
	Registry *naming.Config
	DiaOpts  []grpc.DialOption
	SwgFs    fs.FS
}
