package queue

import "go.uber.org/zap"

/**
以redis为例，一般性理解:
每隔interval拉取step的任务。默认情况下1秒钟拉取1个分数值的任务
假设，同一秒内，添加了56项任务，则：“1秒钟拉取1个分数值的任务”相当于，一次性需要有56个任务要执行
*/

type Options struct {
	//每一次从队列中拉取出来了的间隔（不是个数）
	step int64
	//每隔多久去拉一次队列,单位是秒。一般都是1秒钟。
	interval int64

	logger *zap.SugaredLogger
}

// Option Option
type Option func(o *Options)

func Step(s int64) Option {
	return func(o *Options) {
		o.step = s
	}
}

func Interval(s int64) Option {
	return func(o *Options) {
		o.interval = s
	}
}

func Logger(s *zap.SugaredLogger) Option {
	return func(o *Options) {
		o.logger = s
	}
}

// DefaultOptions default options
func DefaultOptions() Options {
	o := Options{
		step:     1,
		interval: 1,
	}
	return o
}
