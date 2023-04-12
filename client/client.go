package client

import (
	"context"
	"time"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	// DialTimeout the timeout of create connection
	DialTimeout = 5 * time.Second

	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay = 3 * time.Second

	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 1GB is to provide system's throughput.
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize we set it 1GB is to provide system's throughput.
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 4 << 30
)

type Client struct {
	conf    *config.ServerConfig
	org     string
	s       naming.Resolver
	r       naming.Registry
	dialOpt []grpc.DialOption
}

func NewClient(conf *config.ServerConfig, dialOpt ...grpc.DialOption) *Client {

	return &Client{
		conf:    conf,
		org:     conf.Registry.Org,
		s:       naming.NewResolver(conf.Registry),
		r:       naming.NewRegistry(conf.Registry),
		dialOpt: dialOpt,
	}

}

// 没有对外的调用，目前只支持不带验证的
// 默认执行
func (e *Client) Invoke(ctx context.Context, app, method string, in, rs interface{}, h ...string) error {

	//统一独立部署，只有一个target
	target := app

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}
	opts = append(opts, e.dialOpt...)
	conn, err := grpc.Dial(e.r.Scheme()+"://"+target, opts...)
	if err != nil {
		return err
	}
	//handler := cases.Title(language.English).String(app) + "Handler"
	handler := "Handler"
	if len(h) > 0 {
		handler = h[0] + "Handler"
	}
	//[org].[app].[Handler].[method]
	//fmt.Println("/" + e.org + "." + app + "." + handler + "/" + method)
	return conn.Invoke(ctx, "/"+e.org+"."+app+"."+handler+"/"+method, in, rs)
}
