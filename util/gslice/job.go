package gslice

import "time"

// 分批执行一个slice的任务
func RangeJob[T Number](ids []T, step int, job func([]T), duration ...time.Duration) {
	l := len(ids)
	for i := 0; i < l; {
		h := i + step
		if h > l {
			h = l
		}
		if len(duration) > 0 {
			time.Sleep(duration[0])
		}
		go job(ids[i:h])
		i = i + step
	}
}
