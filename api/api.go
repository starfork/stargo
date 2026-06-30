package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Api struct {
	conf   *Config
	conn   *grpc.ClientConn
	ctx    context.Context
	rmux   *runtime.ServeMux
	mux    *http.ServeMux
	server *http.Server
}

func NewApi(conf *Config) (*Api, error) {
	ctx := context.Background()

	var conn *grpc.ClientConn
	if conf.Registry != nil {
		r, err := naming.NewResolver(conf.Registry.Scheme, conf.Registry)
		if err != nil {
			return nil, err
		}
		if len(conf.DiaOpts) == 0 {
			conf.DiaOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		}

		conn, err = client.New(ctx, r, logger.DefaultLogger).NewClient(conf.App, conf.DiaOpts...)
		if err != nil {
			return nil, err
		}
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
	}, nil
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

// DefaultHandlerWrapper is a pass-through wrapper (used as fallback when no custom wrapper is set)
var DefaultHandlerWrapper = func(h http.Handler) http.Handler { return h }

func (e *Api) MuxHandler() {

}
func (e *Api) Run() error {

	if len(e.conf.MuxHandler) > 0 {
		for r, f := range e.conf.MuxHandler {
			e.mux.HandleFunc(r, f)
		}
	}
	e.mux.Handle("/metrics", promhttp.Handler())
	e.mux.Handle("/", e.rmux)

	//e.WrapperSwagger(e.mux)
	// start a standard HTTP server with the router

	var handler http.Handler = e.mux

	// Apply CORS wrapper if configured
	if e.conf.CORS != nil {
		handler = CORSWrapper(*e.conf.CORS)(handler)
	}

	// Apply custom wrapper if configured
	wrapper := DefaultHandlerWrapper
	if e.conf.Wrapper != nil {
		wrapper = e.conf.Wrapper
	}
	handler = wrapper(handler)

	e.server = &http.Server{
		Addr:              e.conf.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Println("start listen " + e.conf.Port)

	// Serve with TLS if configured
	if e.conf.CertFile != "" && e.conf.KeyFile != "" {
		return e.server.ListenAndServeTLS(e.conf.CertFile, e.conf.KeyFile)
	}
	return e.server.ListenAndServe()
}

func (e *Api) Stop(ctx context.Context) error {
	if e.server != nil {
		return e.server.Shutdown(ctx)
	}
	return nil
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
