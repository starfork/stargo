package queue

import (
	"sync"
	"time"
)

type Queue struct {
	store    Store
	handlers *sync.Map
	opts     Options

	//interval   time.Duration //间隔时段
	//workers int           //最大处理数
}

func New(store Store, opts ...Option) *Queue {
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
func (q *Queue) Register(tag string, h Handler) {
	q.handlers.Store(tag, h)
}

// Load 任务内容
func (q *Queue) Load(handlers map[string]Handler) *Queue {
	for k, v := range handlers {
		q.handlers.Store(k, v)
	}
	return q
}

func (q *Queue) Pop(t *Task) error {
	return q.store.Pop(t)
}
func (q *Queue) Push(t *Task) error {
	return q.store.Push(t)
}

// 一般意义上来说，重复添加一个队列，即表示一个更新
//
//	func (q *Queue) Update(t *Task) error {
//		return q.store.Update(t)
//	}
func (q *Queue) run() {
	t := time.NewTicker(time.Second * time.Duration(q.opts.interval)) //TODO，传入配置，interval
	defer t.Stop()
	for {
		<-t.C
		q.exec()
	}
}
func (q *Queue) exec() {
	rs, err := q.store.FetchJob(q.opts.step)

	if err != nil {
		q.log("任务队列获取失败%+v", err)
	}

	for _, v := range rs {
		task, err := q.store.ReadTask(v)
		if err != nil {
			continue
		}
		q.store.Pop(task)
		hander, ok := q.handlers.Load(task.Tag)
		if !ok {
			q.log("任务获取失败%+v", task.Tag)
			//log
			continue
		}
		//开启goroutine ==> workers
		go func() {
			q.log("sart task ")
			//执行成功则删除任务，否则如果设置了
			if err := hander.(Handler)(task); err != nil {
				q.log("执行失败%+v", err)
				//如果有循环条件设置。则循环加入
				task.Retry++
				if task.TTL > 0 && task.Retry <= task.RetryMax {
					task.Delay = task.TTL
					q.store.Update(task)
				} else {
					q.store.Pop(task)
				}

			} else {
				q.store.Pop(task)
			}
		}()
	}
}

func (q *Queue) log(template string, args ...interface{}) {
	if q.opts.logger != nil {
		q.opts.logger.Debugf(template, args...)
	}
}
