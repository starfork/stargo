package client

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/starfork/stargo/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
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
	org      string
	rpcConfs map[string]*config.Server

	dialOpt map[string][]grpc.DialOption
	conns   map[string]grpc.ClientConnInterface

	mu sync.Mutex
}

func New(conf *config.Config, dialOpt ...map[string][]grpc.DialOption) *Client {

	c := &Client{
		//conf:  conf,
		org: conf.Org,
		//s:     naming.NewResolver(conf.Registry),
		//r:     naming.NewRegistry(conf.Registry),
		conns:    make(map[string]grpc.ClientConnInterface),
		rpcConfs: make(map[string]*config.Server),
		//dialOpt: dialOpt,

	}
	if len(dialOpt) > 0 {
		c.dialOpt = dialOpt[0]
	}
	for k, v := range conf.RpcServer {
		c.rpcConfs[k] = v
	}

	return c

}

func DefaultOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		//在stargo中，似乎可以使用“永久”连接，否则需要设置配置
		//避免  ERROR: [transport] Client received GoAway with error code ENHANCE_YOUR_CALM 这个报错
		// grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// 	Time:                KeepAliveTime,
		// 	Timeout:             KeepAliveTimeout,
		// 	PermitWithoutStream: true,
		// }),
	}
	return opts
}

// // 获取一个连接
func (e *Client) Connection(ctx context.Context, app string, appendOpts ...[]grpc.DialOption) (grpc.ClientConnInterface, error) {

	endpoint, err := e.endpoint(app)
	if err != nil {
		return nil, err
	}

	conn, ok := e.conns[app]
	if !ok {
		opts := DefaultOptions()
		if opt, ok := e.dialOpt[app]; ok {
			opts = append(opts, opt...)
		}
		//扩展的配置
		if len(appendOpts) > 0 {
			opts = append(opts, appendOpts[0]...)
		}
		conn1, err := grpc.DialContext(ctx, endpoint, opts...)
		if err != nil {
			return nil, err
		}
		//defer conn.Close()
		e.mu.Lock()
		defer e.mu.Unlock()
		e.conns[app] = conn1

		defer func() {
			if err != nil {
				if cerr := conn1.Close(); cerr != nil {
					grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
				}
				return
			}
			go func() {
				<-ctx.Done()
				if cerr := conn1.Close(); cerr != nil {
					grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
				}
			}()
		}()
	}

	return conn, nil
}
func (e *Client) endpoint(app string) (string, error) {
	conf, ok := e.rpcConfs[app]
	if !ok {
		return "", errors.New("unknow app")
	}
	return conf.Name + "://" + app, nil
}

// depreciated
// 此方法不好，虽然看起来封装了，但是当method改名或者更换之后，编辑器无法识别
// 没有对外的调用，目前只支持不带验证的
// 默认执行
func (e *Client) Invoke(ctx context.Context, app, method string, in, rs interface{}, h ...string) error {
	endpoint, err := e.endpoint(app)
	if err != nil {
		return err
	}

	opts := DefaultOptions()
	if opt, ok := e.dialOpt[app]; ok {
		opts = append(opts, opt...)
	}
	//fmt.Println(e.r.Scheme() + "://" + target)
	//似乎不用每次都Dial
	conn, err := grpc.Dial(endpoint, opts...)

	if err != nil {
		return err
	}

	defer conn.Close()

	//handler := cases.Title(language.English).String(app) + "Handler"
	handler := "Handler"
	if len(h) > 0 {
		handler = h[0] + "Handler"
	}
	//[org].[app].[Handler].[method]
	rpcMethod := "/" + e.org + "." + app + "." + handler + "/" + method

	return conn.Invoke(ctx, rpcMethod, in, rs)
}

// func (e *Client) Close(app ...string) {
// 	if len(app) > 0 {
// 		if conn, ok := e.conns[app[0]]; ok {
// 			conn.Close()
// 		}
// 	} else {
// 		for _, v := range e.conns {
// 			if v != nil {
// 				v.Close()
// 			}
// 		}
// 	}
// }
