package queue

import (
	"fmt"
	"time"

	driver "github.com/starfork/stargo/queue/store/redis"

	jsoniter "github.com/json-iterator/go"
)

type Task struct {
	Key   string                 // system_name/job_id
	Delay int64                  // 延迟时间
	Cycle bool                   // 是否周期循环。
	Tag   string                 // Tag匹配Handler，无Tag的Task将不会被执行
	Args  map[string]interface{} // 任务参数
}

type Handler func(*Task) error

func (t *Task) MarshalJson() string {
	task, _ := jsoniter.Marshal(t)
	return string(task)
}

func UnmarshalJson(j string) (*Task, error) {
	var task Task
	err := jsoniter.Unmarshal([]byte(j), &task)
	return &task, err
}

type Store struct {
	store StoreInterface
}
type StoreInterface interface {
	AddJob(string, string, string, float64) error
	FetchJob() []string                       //拉取任务队列
	FetchTask(string) string                  //获取单个执行任务。
	RemoveTask(name string, key string) error //删除任务
}

func NewStore(store *driver.Storage) *Store {
	return &Store{store: store}
}

func (e *Store) FetchJob() []string {
	return e.store.FetchJob()
}
func (e *Store) AddJob(t *Task) error {
	str := t.MarshalJson()
	delay := time.Now().Unix() + t.Delay
	fmt.Println("add job time" + time.Now().Format("2006-01-02 15:04:05"))
	return e.store.AddJob(t.Tag, t.Key, str, float64(delay))
}

func (e *Store) FetchTask(name string) *Task {
	rs := e.store.FetchTask(name)
	task, _ := UnmarshalJson(rs)
	return task
}
func (e *Store) RemoveTask(t *Task) error {
	return e.store.RemoveTask(t.Tag, t.Key)

}
