package limiter

import (
	"log"
	"sync"
	"testing"
	"time"

	"go.uber.org/ratelimit"
)

var testTks = []string{"aaaa", "bbbb", "cccc", "dddd", "eeee"}
var testTkStore sync.Map

func TestLimiter(t *testing.T) {
	for _, v := range testTks {
		rl := NewLimiter(Policy{
			Tk:  v,
			Num: 1,
		})
		go func(rl *Limiter, v string) {
			for i := 0; i < 5; i++ {
				if !rl.Allow() {
					log.Print("forbid", v)

				} else {
					log.Print("--", v)
				}
			}
		}(rl, v)
	}
	time.Sleep(time.Second * 5) //跟测试时间一样长看结果
}

func TestUberLimiter(t *testing.T) {

	for i, v := range testTks {
		go doTask(i, v)
	}

	time.Sleep(time.Second * 5) //跟测试时间一样长看结果

}

func doTask(i int, v string) {
	limiter, _ := testTkStore.LoadOrStore(v, ratelimit.New(i+3, ratelimit.WithoutSlack))
	rl := limiter.(ratelimit.Limiter)
	prev := time.Now()
	for i := 0; i < 5; i++ {
		now := rl.Take()
		log.Print(v, "--", i, now.Sub(prev))
		prev = now
	}

}
