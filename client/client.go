package client

import (
	"context"
	"fmt"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	fmt.Println("/" + e.org + "." + app + "." + handler + "/" + method)
	return conn.Invoke(ctx, "/"+e.org+"."+app+"."+handler+"/"+method, in, rs)
}
