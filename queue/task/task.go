package task

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type Task struct {
	Key   string         // system_name/job_id
	Tag   string         // Tag匹配Handler，无Tag的Task将不会被执行
	Args  map[string]any // 任务参数
	Delay int64          // 延迟时间

	TTL      int64 //如果执行失败，下一次重复的时间
	Retry    int   //已经重试的次数
	RetryMax int   //最大重试次数
}

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
