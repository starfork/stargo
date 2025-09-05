package api

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
)

type Config struct {
	App      string //
	Port     string
	Registry *naming.Config
	DiaOpts  []grpc.DialOption
	Wrapper  func(http.Handler) http.Handler

	SMOpts []runtime.ServeMuxOption

	MuxHandler map[string]func(w http.ResponseWriter, r *http.Request)

	Enc    bool   //是否开启API数据混淆/加密
	EncKey string //加密key
}
