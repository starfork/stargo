package store

import "github.com/starfork/stargo/queue/task"

//	type Store struct {
//		store StoreInterface
//	}
type Store interface {
	Push(t *task.Task) error   //添加任务
	Pop(t *task.Task) error    //剔除任务
	Update(t *task.Task) error //更新任务--重复调用Push其实就是update了。感觉这个有点多余
	//获取单个执行任务。
	ReadTask(key string) (*task.Task, error)

	//拉取所有任务队列.返回任务名称
	FetchJob(step int64) ([]string, error)
}
