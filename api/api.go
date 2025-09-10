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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Api struct {
	conf *Config
	conn *grpc.ClientConn
	ctx  context.Context
	rmux *runtime.ServeMux
	mux  *http.ServeMux
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

	var conn *grpc.ClientConn
	//手动注册的模式
	if conf.Registry != nil {
		r, err := etcd.NewResolver(conf.Registry)
		E(err)
		if len(conf.DiaOpts) == 0 {
			conf.DiaOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		}

		conn, err = client.New(ctx, r, logger.DefaultLogger).NewClient(conf.App, conf.DiaOpts...)
		E(err)
	}

	if len(conf.SMOpts) == 0 {
		conf.SMOpts = append(conf.SMOpts, DefaultMarshaler)
	}
	rmux := runtime.NewServeMux(conf.SMOpts...)

	mux := http.NewServeMux()
	return &Api{
		conf: conf,
		ctx:  ctx,
		conn: conn,
		rmux: rmux,
		mux:  mux,
	}
}

var DefaultMarshalerOption = runtime.JSONPb{
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames:   true, // 使用 proto 定义里的名字（snake_case）
		EmitUnpopulated: true, // 输出默认值字段
	},
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true,
	},
}

var DefaultMarshaler = runtime.WithMarshalerOption(runtime.MIMEWildcard, &DefaultMarshalerOption)

func (e *Api) MuxHandler() {

}
func (e *Api) Run() {

	if len(e.conf.MuxHandler) > 0 {
		for r, f := range e.conf.MuxHandler {
			e.mux.HandleFunc(r, f)
		}
	}
	e.mux.Handle("/", e.rmux)

	//e.WrapperSwagger(e.mux)
	// start a standard HTTP server with the router
	log.Println("start listen " + e.conf.Port)
	if err := http.ListenAndServe(e.conf.Port, e.conf.Wrapper(e.mux)); err != nil {
		log.Fatal(err)
	}
}

func (e *Api) SetConn(conn *grpc.ClientConn) {
	e.conn = conn
}

func (e *Api) Conn() *grpc.ClientConn {
	return e.conn
}
func (e *Api) Ctx() context.Context {
	return e.ctx
}
func (e *Api) SetCtx(ctx context.Context) {
	e.ctx = ctx
}
func (e *Api) Rmux() *runtime.ServeMux {
	return e.rmux
}
