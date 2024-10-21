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
	"github.com/starfork/stargo/naming"
	"go.etcd.io/etcd/client/v3/naming/resolver"
)

type Resolver struct {
}

// NewResolver creates a resolver builder.
func NewResolver(conf *naming.Config) (naming.Resolver, error) {
	cli, err := newClient(conf)
	if err != nil {
		return nil, err
	}

	return resolver.NewBuilder(cli)
}
