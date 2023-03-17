package limiter

import (
	"net/http"
	"strings"
	"sync"

	"github.com/starfork/stargo/util/request"
	"go.uber.org/ratelimit"
)

var store sync.Map

type Limiter struct {
	p Policy
	//store map[string]*ratelimit.Limiter
	//visits  uint64

}
type Policy struct {
	Type string
	Tk   string
	Num  int //per second
	Req  *http.Request
}

func NewLimiter(p Policy) *Limiter {

	return &Limiter{
		p: p,
	}
}

func (e *Limiter) Allow() bool {
	lr := e.getLimier()

	return lr.Take().Unix() > 0
}

func (e *Limiter) getLimier() ratelimit.Limiter {
	tk, err := e.getToken()
	if err != nil {
		panic("token get error") //不满足limiter生成条件
	}
	limiter, _ := store.LoadOrStore(tk, ratelimit.New(e.p.Num, ratelimit.WithoutSlack))
	return limiter.(ratelimit.Limiter)
}

func (e *Limiter) getToken() (string, error) {

	if e.p.Tk != "" {
		return e.p.Tk, nil
	}

	var tk string
	var err error
	t := strings.ToLower(e.p.Type)
	if t == "both" || t == "" {
		if tk, err = request.GetToken(e.p.Req); err != nil {
			tk, err = request.GetIp(e.p.Req)
		}

	} else if t == "ip" {
		tk, err = request.GetIp(e.p.Req)
	} else if t == "token" {
		tk, err = request.GetToken(e.p.Req)
	}

	return tk, err

}
