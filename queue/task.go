package queue

import jsoniter "github.com/json-iterator/go"

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
