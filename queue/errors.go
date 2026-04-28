package queue

var (
	ErrFailGetJob  = "job  get   %+v \r\n"
	ErrFailGetTask = "task get %+v %s \r\n"
	ErrTaskExec    = "task exec  %s %s %s \r\n"
	TaskUpdate     = "task update %s\r\n"
)
