package circuitbreaker

import "time"

type Config struct {
	FailureThreshold float64
	SuccessThreshold int64
	OpenTimeout      time.Duration
	HalfOpenMaxReqs  int64
	WindowSize       time.Duration
	OnStateChange    func(from, to State)
}

func DefaultConfig() *Config {
	return &Config{
		FailureThreshold: 0.5,
		SuccessThreshold: 3,
		OpenTimeout:      30 * time.Second,
		HalfOpenMaxReqs:  5,
		WindowSize:       60 * time.Second,
	}
}

type Option func(*Config)

func WithFailureThreshold(v float64) Option {
	return func(c *Config) {
		c.FailureThreshold = v
	}
}

func WithSuccessThreshold(v int64) Option {
	return func(c *Config) {
		c.SuccessThreshold = v
	}
}

func WithOpenTimeout(v time.Duration) Option {
	return func(c *Config) {
		c.OpenTimeout = v
	}
}

func WithHalfOpenMaxReqs(v int64) Option {
	return func(c *Config) {
		c.HalfOpenMaxReqs = v
	}
}

func WithWindowSize(v time.Duration) Option {
	return func(c *Config) {
		c.WindowSize = v
	}
}

func WithOnStateChange(fn func(from, to State)) Option {
	return func(c *Config) {
		c.OnStateChange = fn
	}
}
