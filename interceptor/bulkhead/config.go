package bulkhead

import "time"

type Config struct {
	MaxConcurrent int64
	WaitTimeout   time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrent: 100,
		WaitTimeout:   0,
	}
}

type Option func(*Config)

func WithMaxConcurrent(v int64) Option {
	return func(c *Config) {
		c.MaxConcurrent = v
	}
}

func WithWaitTimeout(v time.Duration) Option {
	return func(c *Config) {
		c.WaitTimeout = v
	}
}
