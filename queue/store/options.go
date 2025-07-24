package store

import "github.com/starfork/stargo/logger"

/*
以redis为例，一般性理解:
每隔interval拉取step的任务。默认情况下1秒钟拉取1个分数值的任务
假设，同一秒内，添加了56项任务，则：“1秒钟拉取1个分数值的任务”相当于，一次性需要有56个任务要执行
*/

type Options struct {
	Name   string
	Logger logger.Logger
}

// Option Option
type Option func(o *Options)

func WithName(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func WithLogger(s logger.Logger) Option {
	return func(o *Options) {
		o.Logger = s
	}
}

// DefaultOptions default options
func DefaultOptions() Options {
	o := Options{
		Logger: logger.DefaultLogger,
	}
	return o
}
