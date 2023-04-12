package client

import (
	"context"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conf *config.ServerConfig
	org  string
	s    naming.Resolver
	r    naming.Registry
}

func NewClient(conf *config.ServerConfig) *Client {
	return &Client{
		conf: conf,
		org:  conf.Registry.Org,
		s:    naming.NewResolver(conf.Registry),
		r:    naming.NewRegistry(conf.Registry),
	}
}

// 没有对外的调用，目前只支持不带验证的
// 默认执行
func (e *Client) Invoke(ctx context.Context, app, method string, in, rs interface{}, h ...string) error {

	//统一独立部署，只有一个target
	target := app

	conn, err := grpc.Dial(e.r.Scheme()+"://"+target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithResolvers(e.s))
	if err != nil {
		return err
	}

	//handler := cases.Title(language.English).String(app) + "Handler"
	handler := "Handler"
	if len(h) > 0 {
		handler = h[0] + handler
	}
	//[org].[app].[Handler].[method]
	return conn.Invoke(ctx, "/"+e.org+"."+app+"."+handler+"/"+method, in, rs)
}
