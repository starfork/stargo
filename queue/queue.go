package queue

import (
	"sync"
	"time"

	"github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/queue/task"
)

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
func (q *Queue) Register(tag string, h task.Handler) {
	q.handlers.Store(tag, h)
}

// Load 任务内容
func (q *Queue) Load(handlers map[string]task.Handler) *Queue {
	for k, v := range handlers {
		q.handlers.Store(k, v)
	}
	return q
}

func (q *Queue) Pop(t *task.Task) error {
	return q.store.Pop(t)
}
func (q *Queue) Push(t *task.Task) error {
	return q.store.Push(t)
}

// 一般意义上来说，重复添加一个队列，即表示一个更新
//
//	func (q *Queue) Update(t *Task) error {
//		return q.store.Update(t)
//	}
func (q *Queue) run() {
	t := time.NewTicker(time.Second * time.Duration(q.opts.interval))
	defer t.Stop()
	for {
		<-t.C
		q.exec()
	}
}
func (q *Queue) exec() {
	rs, err := q.store.FetchJob(q.opts.step)

	if err != nil {
		q.log(ErrFailGetTask, err)
	}

	for _, v := range rs {
		t, err := q.store.ReadTask(v)
		if err != nil {
			continue
		}
		q.store.Pop(t)
		hander, ok := q.handlers.Load(t.Tag)
		if !ok {
			q.log(ErrFailGetTask, t.Tag)
			//log
			continue
		}
		//开启goroutine ==> workers
		go func() {
			q.log("sart task ")
			//执行成功则删除任务，否则如果设置了
			//fmt.Println(task)
			if err := hander.(task.Handler)(t); err != nil {
				q.log(ErrFailGetTask, err)
				//如果有循环条件设置。则循环加入
				t.Retry++
				if t.TTL > 0 && t.Retry <= t.RetryMax {
					t.Delay = t.TTL
					q.store.Update(t)
				} else {
					q.store.Pop(t)
				}

			} else {
				q.store.Pop(t)
			}
		}()
	}
}

func (q *Queue) log(template string, args ...interface{}) {
	if q.opts.logger != nil {
		q.opts.logger.Debugf(template, args...)
	}
}
