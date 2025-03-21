package client

import (
	"context"
	"fmt"
	"time"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
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
	ctx      context.Context
	resolver naming.Resolver
	//rpcConfs map[string]*config.RpcServer

	logger logger.Logger
}

func New(ctx context.Context, resolver naming.Resolver, logger logger.Logger) *Client {

	c := &Client{
		ctx:      ctx,
		resolver: resolver,
		logger:   logger,
	}

	return c

}

func DefaultOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		//grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
	}
	return opts
}

// 获取一个连接
func (e *Client) NewClient(service string, options ...grpc.DialOption) (conn *grpc.ClientConn, err error) {

	opts := DefaultOptions()
	opts = append(opts, grpc.WithResolvers(e.resolver))
	opts = append(opts, options...)
	target := fmt.Sprintf("%s:///%s/%s", e.resolver.Scheme(), e.resolver.Config().Org, service)
	return grpc.NewClient(target, opts...)
}
