package etcd

import (
	"context"
	"time"

	"github.com/starfork/stargo/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func newClient(conf *config.Registry) (cli *clientv3.Client) {
	var err error
	config := clientv3.Config{
		Endpoints:   []string{conf.Addr},
		DialTimeout: 2 * time.Second,
	}
	if cli, err = clientv3.New(config); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := cli.Status(ctx, config.Endpoints[0]); err != nil {
		panic(err)
	}
	return cli
}
