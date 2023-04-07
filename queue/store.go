package queue

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type Task struct {
	Key   string                 // system_name/job_id
	Tag   string                 // Tag匹配Handler，无Tag的Task将不会被执行
	Args  map[string]interface{} // 任务参数
	Delay int64                  // 延迟时间

	TTL      int64 //如果执行失败，下一次重复的时间
	Retry    int   //已经重试的次数
	RetryMax int   //最大重试次数
}

// 任务执行函数。
// 特别注意，如果ttl不是0。但是因为某些原因不需要执行下一次操作了。则需要返回nil
// 比如。需要对某特定条件数据执行某个操作。在这之前自然是会检查这个数据是否存在或者是否能被操作等。
// 如果没有拿到这条数据，则直接返回nil
type Handler func(*Task) error

func (t *Task) MarshalJson() string {
	task, _ := jsoniter.Marshal(t)
	return string(task)
}

func (t *Task) Subkey() string {
	return t.Tag + "." + t.Key
}

// 把key变成sn。很多情况下，都是sn设为key的
func (t *Task) Sn() uint64 {
	sn, _ := strconv.Atoi(t.Key)
	return uint64(sn)
}

func UnmarshalJson(j string) (*Task, error) {
	var task Task
	err := jsoniter.Unmarshal([]byte(j), &task)
	return &task, err
}

//	type Store struct {
//		store StoreInterface
//	}
type Store interface {
	Push(t *Task) error   //添加任务
	Pop(t *Task) error    //剔除任务
	Update(t *Task) error //更新任务--重复调用Push其实就是update了。感觉这个有点多余
	//获取单个执行任务。
	ReadTask(key string) (*Task, error)

	//拉取所有任务队列.返回任务名称
	FetchJob(step int64) ([]string, error)
}
