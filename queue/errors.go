package queue

var (
	ErrFailGetJob  = "任务执行内容获取失败 %+v \r\n"
	ErrFailGetTask = "任务队列获取失败 %+v \r\n"
	ErrTaskExec    = "任务执行失败 %+v \r\n"
	TaskUpdate     = "任务重新加入执行队列 \r\n"
)
