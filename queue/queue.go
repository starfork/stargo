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
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, e.opts.maxThread)
	)
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

		wg.Add(1)
		sem <- struct{}{} // 获取令牌

		go func(t *task.Task, handler task.Handler) {
			defer func() {
				if r := recover(); r != nil {
					e.log("panic in handler %s: %v", t.Subkey(), r)
					e.store.Pop(t) //直接删掉
				}
				<-sem
				wg.Done()
			}()

			//e.log("[task start] %s", t.Subkey())

			if err := handler(t); err != nil {
				e.log(ErrTaskExec+time.Now().String(), err)
				ttl := t.GetTTL(t.Retry)
				t.Retry++
				if ttl > 0 && t.Retry <= t.RetryMax {
					t.Delay = ttl
					if err := e.store.Update(t); err != nil {
						e.log("[task update] %s failed: %v", t.Subkey(), err)
					} else {
						e.log(TaskUpdate)
					}
				} else {
					e.log("[task finished] %s error %+s", t.Subkey(), err.Error())
					e.store.Pop(t)
				}
			} else {
				//e.log("[task finished] %s success", t.Subkey())
				e.store.Pop(t)
			}
		}(t, hander.(task.Handler))
	}

	wg.Wait()
}

func (e *Queue) log(template string, args ...any) {
	if e.opts.logger != nil {
		start := time.Now()
		e.opts.logger.Debugf(start.Format(tformat)+" "+template+" \r\n", args...)
	}
}
