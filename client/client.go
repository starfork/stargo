package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/starfork/stargo/interceptor/timeout"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/util/ustring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	keepalive "google.golang.org/grpc/keepalive"
)

const (
	// DialTimeout the timeout of create connection
	DialTimeout = 3 * time.Second
	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay = 3 * time.Second
	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = 10 * time.Second
	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = 3 * time.Second

	// InitialWindowSize we set it to 16MB for reasonable throughput.
	InitialWindowSize = 1 << 24 // 16MB

	// InitialConnWindowSize we set it to 16MB for reasonable throughput.
	InitialConnWindowSize = 1 << 24 // 16MB

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 64 << 20 // 64MB

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 64 << 20 // 64MB
)

type Client struct {
	ctx      context.Context
	resolver naming.Resolver
	logger   logger.Logger
	conns    map[string]*grpc.ClientConn
	mu       sync.RWMutex
}

func New(ctx context.Context, resolver naming.Resolver, logger logger.Logger) *Client {

	c := &Client{
		ctx:      ctx,
		resolver: resolver,
		logger:   logger,
		conns:    make(map[string]*grpc.ClientConn),
	}

	return c

}

func DefaultOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingConfig": [{"round_robin":{}}],
			"healthCheckConfig": {"serviceName": ""},
			"methodConfig": [{
				"name": [{"service": ""}],
				"waitForReady": true,
				"retryPolicy": {
					"maxAttempts": 3,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2.0,
					"retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
				}
			}]
		}`),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   BackoffMaxDelay,
			},
			MinConnectTimeout: DialTimeout,
		}),
		grpc.WithChainUnaryInterceptor(timeout.UnaryClient()),
		grpc.WithChainStreamInterceptor(timeout.StreamClient()),
	}
	return opts
}

// 获取一个连接（带缓存）
func (e *Client) NewClient(service string, options ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	target := fmt.Sprintf("%s:///%s/%s", e.resolver.Scheme(), ustring.OrString("stargo", e.resolver.Config().Org), service)
	
	// Check cache first
	e.mu.RLock()
	if conn, ok := e.conns[target]; ok {
		e.mu.RUnlock()
		return conn, nil
	}
	e.mu.RUnlock()
	
	// Create new connection
	var opts []grpc.DialOption
	if len(options) == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, options...)
	}

	opts = append(opts, DefaultOptions()...)
	opts = append(opts, grpc.WithResolvers(e.resolver))

	e.logger.Infof("[client] connecting to %s", target)
	conn, err = grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}
	
	// Cache the connection
	e.mu.Lock()
	e.conns[target] = conn
	e.mu.Unlock()
	
	return conn, nil
}

// Close closes all cached connections
func (e *Client) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	for target, conn := range e.conns {
		if err := conn.Close(); err != nil {
			e.logger.Errorf("[client] failed to close connection to %s: %v", target, err)
		}
		delete(e.conns, target)
	}
	return nil
}
