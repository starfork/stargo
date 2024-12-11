package api

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming/etcd"
	"google.golang.org/grpc"
)

type Api struct {
	conf *Config
	conn *grpc.ClientConn
	ctx  context.Context
	rmux *runtime.ServeMux
}

func E(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func NewApi(conf *Config) *Api {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r, err := etcd.NewResolver(conf.Registry)
	E(err)
	conn, err := client.New(ctx, r, logger.DefaultLogger).NewClient(conf.Registry.Org+"/"+conf.App, conf.DiaOpts...)
	E(err)
	rmux := runtime.NewServeMux(conf.SMOpts...)
	return &Api{conf: conf, ctx: ctx, conn: conn, rmux: rmux}
}

func (e *Api) Run() {
	mux := http.NewServeMux()

	mux.Handle("/", e.rmux)
	e.WrapperSwagger(mux)
	// start a standard HTTP server with the router
	log.Println("start listen " + e.conf.Port)
	// if e.conf.Wrapper != nil {
	// 	e.conf.Wrapper(mux)
	// }
	if err := http.ListenAndServe(e.conf.Port, e.conf.Wrapper(mux)); err != nil {
		log.Fatal(err)
	}
}

func (e *Api) Conn() *grpc.ClientConn {
	return e.conn
}
func (e *Api) Ctx() context.Context {
	return e.ctx
}

func (e *Api) Rmux() *runtime.ServeMux {
	return e.rmux
}
