package queue

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Queue struct {
	store    Store
	handlers *sync.Map
	logger   *zap.SugaredLogger
	//interval   time.Duration //间隔时段
	//workers int           //最大处理数
}

func New(store Store, logger ...*zap.SugaredLogger) *Queue {
	q := &Queue{
		store:    store,
		handlers: &sync.Map{},
	}
	if len(logger) > 0 {
		q.logger = logger[0]
	}
	go q.run()
	return q
}

// Push 任务
func (q *Queue) Push(t *Task) {

	str := t.MarshalJson()
	delay := time.Now().Unix() + t.Delay
	//fmt.Println("add job time" + time.Now().Format("2006-01-02 15:04:05"))
	q.store.AddJob(t.Tag, t.Key, str, float64(delay))

}
func (q *Queue) Pop(t Task) {
	q.store.RemoveTask(t.Tag, t.Key)
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
	rs := q.store.FetchJob()
	for _, v := range rs {
		t := q.store.FetchTask(v)
		task, _ := UnmarshalJson(t)
		q.store.RemoveTask(task.Tag, task.Key)
		h, ok := q.handlers.Load(task.Tag)
		if !ok {
			//log
			continue
		}
		//开启goroutine ==> workers
		go func() {
			//hander.(Handler)(task)
			if err := h.(Handler)(task); err != nil {
				//q.store.AddJob(task)
				q.log("执行失败", err)
			}
			//TODO 如果执行错误，那么就重新添加到队列中（根据策略）
		}()
	}
}

func (q *Queue) log(template string, args ...interface{}) {
	if q != nil {
		q.logger.Debugf(template, args...)
	}
}
