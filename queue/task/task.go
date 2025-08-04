package task

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type Task struct {
	Key string // job_id/ 订单号之类的东西
	Tag string // Tag匹配Handler，无Tag的Task将不会被执行

	// 延迟时间。本次执行之后下一次的执行时间。单位秒
	Delay int64

	//如果执行失败，下一次重复的时间.
	//这个可以设置多段不通时间。分别表示第几次retry的时候需要延迟的时间
	//如果retry对应的index超限，则会返回最后一个。
	DelayTTL []int64

	//已经重试的次数
	Retry int

	//最大重试次数
	RetryMax int

	Args map[string]any // 任务参数
}

func (e *Task) MarshalJson() string {
	task, _ := jsoniter.Marshal(e)
	return string(task)
}

func (e *Task) Subkey() string {
	return e.Tag + "." + e.Key
}

// 把key变成sn。很多情况下，都是sn设为key的
func (e *Task) Sn() uint64 {
	sn, _ := strconv.Atoi(e.Key)
	return uint64(sn)
}

func (e *Task) GetTTL(n int) int64 {
	s := e.DelayTTL
	if n >= 0 && n < len(s) {
		return s[n]
	}
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return 0 // 或根据需求返回默认值
}
