package queue

type Store interface {
	AddJob(name, key, value string, interval float64) error
	FetchJob() []string                       //拉取任务队列
	FetchTask(name string) string             //获取单个执行任务。
	RemoveTask(name string, key string) error //删除任务
}
