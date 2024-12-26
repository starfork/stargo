package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/starfork/stargo/naming"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func newClient(conf *naming.Config) (cli *clientv3.Client, err error) {

	config := clientv3.Config{
		Endpoints:   strings.Split(conf.Host, ","),
		DialTimeout: 2 * time.Second,
	}
	if cli, err = clientv3.New(config); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := cli.Status(ctx, config.Endpoints[0]); err != nil {
		return nil, err
	}
	return cli, nil
}
