package request

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type VisitorLimiter struct {
	store  map[string]*Limiter
	policy Policy
	req    *http.Request
	//limiter *rate.Limiter
	mu sync.Mutex
}
type Limiter struct {
	visits  uint64
	limiter *rate.Limiter
}
type Policy struct {
	Name   string
	Ip, Tk string
	R      rate.Limit
	B      int
}

func NewVisitorLimiter() *VisitorLimiter {
	return &VisitorLimiter{
		store: make(map[string]*Limiter),
		policy: Policy{
			Name: "both",
		},
	}
}
func NewLimiter(p Policy) *Limiter {
	return &Limiter{
		visits:  0,
		limiter: rate.NewLimiter(p.R, p.B),
	}
}
func (e *Limiter) Allow() bool {
	//fmt.Printf("%v", e)

	e.visits++
	fmt.Printf("\r limit %d visits", e.visits)
	return e.limiter.Allow()

}
func (e *Limiter) Visits() uint64 {
	return e.visits
}
func (e *VisitorLimiter) Allow() bool {
	var err error
	var limiter *Limiter
	if e.policy.Name == "token" {
		limiter, err = e.NewTokenLimiter()
	} else if e.policy.Name == "ip" {
		limiter, err = e.NewIpLimiter()
	} else {
		limiter, err = e.SetPolicy(Policy{
			Name: "token",
			R:    rate.Every(2 * time.Second),
			B:    20,
		}).NewTokenLimiter()
		if err != nil {
			limiter, err = e.SetPolicy(Policy{
				Name: "ip",
				R:    rate.Every(2 * time.Second),
				B:    30,
			}).NewIpLimiter()
		}
	}
	if err != nil {
		return false //如果token没得，ip没得，直接报频率限制错误
	}
	limiter.visits++
	return limiter.limiter.Allow()
}

func (e *VisitorLimiter) NewIpLimiter() (*Limiter, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var ip string
	var err error
	if e.policy.Ip != "" {
		ip = e.policy.Ip
	} else {
		ip, _, err = net.SplitHostPort(e.req.RemoteAddr)
	}
	if ip == "" || err != nil {
		return nil, errors.New("limiter ip error")
	}
	limiter, exists := e.store[ip]
	//glog.Infof("e.ips%+v", e.ips)
	if !exists {
		limiter := NewLimiter(e.policy)
		//limiter =
		e.store[ip] = limiter
	}
	return limiter, nil
}

func (e *VisitorLimiter) NewTokenLimiter() (*Limiter, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var tk string
	var err error
	if e.policy.Tk != "" {
		tk = e.policy.Tk
	} else {
		tk, err = GetToken(e.req)
	}
	if tk == "" || err != nil {
		return nil, errors.New("limiter tk error")
	}
	limiter, exists := e.store[tk]
	//glog.Infof("e.tks%+v", e.tks)
	if !exists {
		limiter := NewLimiter(e.policy)
		//fmt.Printf("%v", limiter)
		e.store[tk] = limiter
	}
	return limiter, nil
}

func (e *VisitorLimiter) SetRequest(req *http.Request) *VisitorLimiter {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.req = req
	return e
}
func (e *VisitorLimiter) SetPolicy(p Policy) *VisitorLimiter {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.policy = p
	return e
}
