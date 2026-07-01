package queue

import (
	"context"
	"sync"
	"time"

	"github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/queue/task"
)

var tformat = "2006-01-02 15:04:05"

type Queue struct {
	ctx       context.Context
	cancel    context.CancelFunc
	store     store.Store
	handlers  *sync.Map
	opts      Options
	executing bool
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func New(ctx context.Context, store store.Store, opts ...Option) *Queue {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(ctx)
	q := &Queue{
		ctx:      ctx,
		cancel:   cancel,
		store:    store,
		handlers: &sync.Map{},
		opts:     options,
	}
	go q.run()
	return q
}

// Register 执行方法
func (e *Queue) Register(tag string, h task.Handler) {
	e.handlers.Store(tag, h)
}

// Load 任务内容
func (e *Queue) Load(handlers map[string]task.Handler) *Queue {
	for k, v := range handlers {
		e.Register(k, v)
	}
	return e
}

func (e *Queue) Pop(t *task.Task) error {
	return e.store.Pop(t)
}
func (e *Queue) Push(t *task.Task) error {
	return e.store.Push(t)
}

// 一般意义上来说，重复添加一个队列，即表示一个更新
//
//	func (e *Queue) Update(t *Task) error {
//		return e.store.Update(t)
//	}
func (e *Queue) run() {
	t := time.NewTicker(time.Second * time.Duration(e.opts.interval))
	defer t.Stop()
	
	// Reclaim expired jobs ticker (every 10 seconds)
	reclaimTicker := time.NewTicker(10 * time.Second)
	defer reclaimTicker.Stop()
	
	for {
		select {
		case <-e.ctx.Done():
			return
		case <-t.C:
			e.wg.Add(1)
			go func() {
				defer e.wg.Done()
				e.exec()
			}()
		case <-reclaimTicker.C:
			e.wg.Add(1)
			go func() {
				defer e.wg.Done()
				e.reclaim()
			}()
		}
	}
}
func (e *Queue) exec() {
	// Try to acquire lock
	e.mu.Lock()
	if e.executing {
		e.mu.Unlock()
		e.log("exec already in progress, skipping")
		return
	}
	e.executing = true
	e.mu.Unlock()
	
	defer func() {
		e.mu.Lock()
		e.executing = false
		e.mu.Unlock()
	}()
	
	rs, err := e.store.ClaimJob(e.opts.step)
	if err != nil {
		e.log(ErrFailGetJob, err)
		return
	}
	
	if len(rs) == 0 {
		return
	}
	
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, e.opts.maxThread)
	)
	for _, v := range rs {
		t, err := e.store.ReadTask(v)
		if err != nil {
			continue
		}

		handler, ok := e.handlers.Load(t.Tag)
		if !ok {
			e.log(ErrFailGetTask, t.Tag, t.Key)
			continue
		}

		sem <- struct{}{} // 获取令牌

		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					e.log("panic in handler %s: %v", t.Subkey(), r)
					ttl := t.GetTTL(t.Retry)
					t.Retry++
					if t.RetryOnPanic && ttl > 0 && t.Retry <= t.RetryMax {
						t.Delay = ttl
						if err := e.store.Update(t); err != nil {
							e.log("[task panic retry failed] %s: %v", t.Subkey(), err)
							e.store.Pop(t)
						} else {
							e.log("[task panic retry scheduled] %s retry=%d", t.Subkey(), t.Retry)
						}
					} else {
						e.log("[task panic dlq] %s", t.Subkey())
						if dlqErr := e.store.PushDLQ(t); dlqErr != nil {
							e.log("[task panic dlq failed] %s: %v", t.Subkey(), dlqErr)
						}
						e.store.Pop(t)
					}
				}
				<-sem
			}()

			if err := handler.(task.Handler)(t); err != nil {
				e.log(ErrTaskExec, t.Key, t.Tag, err.Error())
				ttl := t.GetTTL(t.Retry)
				t.Retry++
				if ttl > 0 && t.Retry <= t.RetryMax {
					t.Delay = ttl
					if err := e.store.Update(t); err != nil {
						e.log("[task upd err] %s failed: %v", t.Subkey(), err)
					} else {
						e.log("[task upd retry] %s", t.Subkey())
					}
				} else {
					e.log("[task err dlq] %s %s", t.Subkey(), err.Error())
					if dlqErr := e.store.PushDLQ(t); dlqErr != nil {
						e.log("[task dlq failed] %s: %v", t.Subkey(), dlqErr)
					}
					e.store.Pop(t)
				}
			} else {
				e.store.Pop(t)
			}
		})
	}

	wg.Wait()
}

func (e *Queue) log(template string, args ...any) {
	if e.opts.logger != nil {
		start := time.Now()
		e.opts.logger.Debugf(start.Format(tformat)+" "+template+" \r\n", args...)
	}
}

func (e *Queue) reclaim() {
	rs, err := e.store.ReclaimExpired()
	if err != nil {
		e.log("reclaim expired jobs error: %v", err)
		return
	}
	if len(rs) > 0 {
		e.log("reclaimed %d expired jobs", len(rs))
	}
}

func (e *Queue) ListDLQ(offset, limit int64) ([]*task.Task, error) {
	keys, err := e.store.FetchDLQ(offset, limit)
	if err != nil {
		return nil, err
	}
	var tasks []*task.Task
	for _, key := range keys {
		if t, err := e.store.ReadDLQTask(key); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (e *Queue) ReplayFromDLQ(taskKey string) error {
	t, err := e.store.ReadDLQTask(taskKey)
	if err != nil {
		return err
	}
	return e.store.ReplayFromDLQ(t)
}

func (e *Queue) DeleteFromDLQ(key string) error {
	return e.store.PopDLQ(key)
}

func (e *Queue) Stop() {
	e.cancel()
	e.wg.Wait()
}
