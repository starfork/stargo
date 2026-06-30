package etcd

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const WeightedRoundRobin = "weighted_round_robin_xds"

type addrInfo struct {
	sc     balancer.SubConn
	weight int64
}

type wrrPicker struct {
	addrs []*addrInfo
	next  atomic.Int64
	mu    sync.Mutex
}

func init() {
	balancer.Register(NewWRRBuilder())
}

func NewWRRBuilder() balancer.Builder {
	return base.NewBalancerBuilder(WeightedRoundRobin, &wrrPickerBuilder{}, base.Config{HealthCheck: true})
}

type wrrPickerBuilder struct{}

func (*wrrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var addrs []*addrInfo
	for sc, sci := range info.ReadySCs {
		weight := int64(100)
		if battr := sci.Address.BalancerAttributes; battr != nil {
			if w := battr.Value("weight"); w != nil {
				switch wv := w.(type) {
				case int64:
					weight = wv
				case float64:
					weight = int64(wv)
				}
			}
		}
		addrs = append(addrs, &addrInfo{
			sc:     sc,
			weight: weight,
		})
	}

	return &wrrPicker{
		addrs: addrs,
	}
}

func (p *wrrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.addrs) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var totalWeight int64
	for _, ai := range p.addrs {
		totalWeight += ai.weight
	}

	if totalWeight <= 0 {
		idx := p.next.Add(1) % int64(len(p.addrs))
		return balancer.PickResult{SubConn: p.addrs[idx].sc}, nil
	}

	pos := int(p.next.Add(1) % totalWeight)
	for _, ai := range p.addrs {
		if pos < int(ai.weight) {
			return balancer.PickResult{SubConn: ai.sc}, nil
		}
		pos -= int(ai.weight)
	}

	return balancer.PickResult{SubConn: p.addrs[0].sc}, nil
}

func parseMetadata(metadata any) (version string, weight int64) {
	weight = 100
	switch v := metadata.(type) {
	case json.RawMessage:
		var m serviceMetadata
		if err := json.Unmarshal(v, &m); err == nil {
			version = m.Version
			if m.Weight > 0 {
				weight = m.Weight
			}
		}
	case string:
		var m serviceMetadata
		if err := json.Unmarshal([]byte(v), &m); err == nil {
			version = m.Version
			if m.Weight > 0 {
				weight = m.Weight
			}
		}
	case []byte:
		var m serviceMetadata
		if err := json.Unmarshal(v, &m); err == nil {
			version = m.Version
			if m.Weight > 0 {
				weight = m.Weight
			}
		}
	}
	return
}
