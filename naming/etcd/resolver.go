// Copyright 2021 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"sync"

	"github.com/starfork/stargo/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Resolver struct {
	c *clientv3.Client
	//target string
	cc     resolver.ClientConn
	wch    endpoints.WatchChannel
	ctx    context.Context
	cancel context.CancelFunc

	conf *config.Registry
	wg   sync.WaitGroup
}

// type builder struct {
// 	c *clientv3.Client
// }

// NewResolver creates a resolver builder.
func NewResolver(conf *config.Registry) resolver.Builder {
	client := newClient(conf)
	r := &Resolver{c: client, conf: conf}

	resolver.Register(r)
	return r
}
func (e *Resolver) key(name string) string {
	return e.conf.Org + "/" + name
}

func (e *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// Refer to https://github.com/grpc/grpc-go/blob/16d3df80f029f57cff5458f1d6da6aedbc23545d/clientconn.go#L1587-L1611
	key := e.key(target.URL.Host)

	// r := &Resolver{
	// 	//c:      b.c,
	// 	target: key,
	// 	cc:     cc,
	// }
	e.cc = cc
	e.ctx, e.cancel = context.WithCancel(context.Background())

	em, err := endpoints.NewManager(e.c, key)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "resolver: failed to new endpoint manager: %s", err)
	}
	e.wch, err = em.NewWatchChannel(e.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolver: failed to new watch channer: %s", err)
	}

	e.wg.Add(1)
	go e.watch()
	return e, nil
}

func (b *Resolver) Scheme() string {
	return "etcd"
}

func (e *Resolver) watch() {
	defer e.wg.Done()

	allUps := make(map[string]*endpoints.Update)
	for {
		select {
		case <-e.ctx.Done():
			return
		case ups, ok := <-e.wch:
			if !ok {
				return
			}

			for _, up := range ups {
				switch up.Op {
				case endpoints.Add:
					allUps[up.Key] = up
				case endpoints.Delete:
					delete(allUps, up.Key)
				}
			}

			addrs := convertToGRPCAddress(allUps)
			e.cc.UpdateState(resolver.State{Addresses: addrs})
		}
	}
}

func convertToGRPCAddress(ups map[string]*endpoints.Update) []resolver.Address {
	var addrs []resolver.Address
	for _, up := range ups {
		addr := resolver.Address{
			Addr:     up.Endpoint.Addr,
			Metadata: up.Endpoint.Metadata,
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {
	r.cancel()
	r.wg.Wait()
}
