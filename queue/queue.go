package queue

import (
	"sync"
	"time"
)

type Queue struct {
	storage  *Store
	handlers *sync.Map
	//logger   *zap.SugaredLogger
	//interval   time.Duration //间隔时段
	//workers int           //最大处理数
}

func New(storage *Store) *Queue {
	q := &Queue{
		storage:  storage,
		handlers: &sync.Map{},
	}
	go q.run()
	return q
}

// Push 任务
func (q *Queue) Push(t *Task) {
	q.storage.AddJob(t)

}
func (q *Queue) Pop(t Task) {
	q.storage.RemoveTask(&t)
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

func (q *Queue) run() {
	t := time.NewTicker(time.Second * 1) //TODO，传入配置，interval
	defer t.Stop()
	for {
		<-t.C
		q.exec()
	}
}
func (q *Queue) exec() {
	rs := q.storage.FetchJob()
	for _, v := range rs {
		task := q.storage.FetchTask(v)
		q.storage.RemoveTask(task)
		hander, ok := q.handlers.Load(task.Tag)
		if !ok {
			//log
			continue
		}
		//开启goroutine ==> workers
		go func() {
			//hander.(Handler)(task)

			if err := hander.(Handler)(task); err != nil {
				//q.storage.AddJob(task)
				//q.logger.Debug("执行失败", err)
			}
			//TODO 如果执行错误，那么就重新添加到队列中（根据策略）
		}()
	}
}
