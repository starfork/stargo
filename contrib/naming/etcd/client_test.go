package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDial(t *testing.T) {
	conf := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	}
	cli, err := clientv3.New(conf)
	fmt.Println(err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rs, err := cli.Status(ctx, conf.Endpoints[0])
	fmt.Println(rs)
	fmt.Println(err)
	//defer cli.Close()
}
