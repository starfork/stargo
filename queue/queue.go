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

func (q *Queue) run() {
	t := time.NewTicker(time.Second * time.Duration(q.opts.interval)) //TODO，传入配置，interval
	defer t.Stop()
	for {
		<-t.C
		q.exec()
	}
}
func (q *Queue) exec() {
	rs, _ := q.store.FetchJob(q.opts.step)

	for _, v := range rs {
		task, err := q.store.ReadTask(v)
		if err != nil {
			continue
		}
		q.store.Pop(task)
		hander, ok := q.handlers.Load(task.Tag)
		if !ok {
			//log
			continue
		}
		//开启goroutine ==> workers
		go func() {
			//hander.(Handler)(task)
			if err := hander.(Handler)(task); err != nil {
				q.opts.logger.Debugf("执行失败%+v", err)
			}
			// if err := hander.(Handler)(task); err != nil {
			// 	//q.storage.AddJob(task)
			// 	//q.logger.Debug("执行失败", err)
			// }
			//TODO 如果执行错误，那么就重新添加到队列中（根据策略）
		}()
	}
}
