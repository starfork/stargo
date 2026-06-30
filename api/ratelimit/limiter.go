package limiter

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/starfork/stargo/util/request"
	"go.uber.org/ratelimit"
)

type limiterEntry struct {
	limiter  ratelimit.Limiter
	lastSeen atomic.Int64
}

var (
	store     sync.Map
	cleanOnce sync.Once
	stopChan  = make(chan struct{})
)

func init() {
	go cleanup()
}

func cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			now := time.Now().UnixNano()
			store.Range(func(key, value any) bool {
				e := value.(*limiterEntry)
				if now-e.lastSeen.Load() > int64(10*time.Minute) {
					store.Delete(key)
				}
				return true
			})
		}
	}
}

type Limiter struct {
	p Policy
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
	lr, err := e.getLimier()
	if err != nil {
		return false
	}
	return lr.Take().Unix() > 0
}

func (l *Limiter) getLimier() (ratelimit.Limiter, error) {
	tk, err := l.getToken()
	if err != nil {
		return nil, fmt.Errorf("limiter token: %w", err)
	}
	
	entry, _ := store.LoadOrStore(tk, &limiterEntry{
		limiter: ratelimit.New(l.p.Num, ratelimit.WithoutSlack),
	})
	le := entry.(*limiterEntry)
	le.lastSeen.Store(time.Now().UnixNano())
	return le.limiter, nil
}

func (e *Limiter) getToken() (tk string, err error) {

	if e.p.Tk != "" {
		return e.p.Tk, nil
	}

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

	return

}
