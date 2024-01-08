package task

import jsoniter "github.com/json-iterator/go"

// 任务执行函数。
// 特别注意，如果ttl不是0。但是因为某些原因不需要执行下一次操作了。则需要返回nil
// 比如。需要对某特定条件数据执行某个操作。在这之前自然是会检查这个数据是否存在或者是否能被操作等。
// 如果没有拿到这条数据，则直接返回nil
type Handler func(*Task) error

func UnmarshalJson(j string) (*Task, error) {
	var task Task
	err := jsoniter.Unmarshal([]byte(j), &task)
	return &task, err
}
