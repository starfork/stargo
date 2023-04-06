package queue

import jsoniter "github.com/json-iterator/go"

type Store interface {
	Push(t *Task) error   //添加任务
	Pop(t *Task) error    //剔除任务
	Update(t *Task) error //更新任务
	//获取单个执行任务。
	ReadTask(key string) (*Task, error)

	//拉取所有任务队列.返回任务名称
	FetchJob(step int64) ([]string, error)
}

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
