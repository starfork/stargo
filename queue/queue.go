package queue

import (
	"sync"
	"time"

	"github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/queue/task"
)

var tformat = "2006-01-02 15:04:05"

type Queue struct {
	store    store.Store
	handlers *sync.Map
	opts     Options
	//interval   time.Duration //间隔时段
	//workers int           //最大处理数
}

func New(store store.Store, opts ...Option) *Queue {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	q := &Queue{
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
	for {
		<-t.C
		e.exec()
	}
}
func (e *Queue) exec() {
	rs, err := e.store.FetchJob(e.opts.step)
	if err != nil {
		e.log(ErrFailGetJob, err)
	}

	// 限制最大并发数
	sem := make(chan struct{}, e.opts.maxThread)
	var wg sync.WaitGroup

	for _, v := range rs {
		t, err := e.store.ReadTask(v)
		if err != nil {
			continue
		}
		e.store.Pop(t)

		hander, ok := e.handlers.Load(t.Tag)
		if !ok {
			e.log(ErrFailGetTask, t.Tag)
			continue
		}

		sem <- struct{}{} // 占用一个位置
		wg.Add(1)

		go func(t *task.Task, handler task.Handler) {
			defer func() {
				<-sem // 释放位置
				wg.Done()
			}()

			e.log("start task  %s  at %s \r\n", t.Subkey(), time.Now().Format(tformat))

			if err := handler(t); err != nil {
				e.log(ErrTaskExec, err)
				ttl := t.GetTTL(t.Retry)
				t.Retry++
				if ttl > 0 && t.Retry <= t.RetryMax {
					t.Delay = ttl
					e.store.Update(t)
					e.log(TaskUpdate)
				} else {
					e.log("finished task %s error %+s at %s \r\n", t.Subkey(), err.Error(), time.Now().Format(tformat))
					e.store.Pop(t)
				}
			} else {
				e.log("finished task %s success  at %s \r\n", t.Subkey(), time.Now().Format(tformat))
				e.store.Pop(t)
			}
		}(t, hander.(task.Handler))
	}

	wg.Wait()
}

func (e *Queue) log(template string, args ...any) {
	if e.opts.logger != nil {
		e.opts.logger.Debugf(template, args...)
	}
}
