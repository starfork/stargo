package logger

import "fmt"

type Logger interface {
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Infof(format string, v ...any)

	String() string
	Options() Options
}

var DefaultLogger Logger = newDefaultLogger()

var loggerFactories = make(map[string]func(*Config) (Logger, error))

func Register(name string, factory func(*Config) (Logger, error)) {
	loggerFactories[name] = factory
}

func NewLogger(name string, conf *Config) (Logger, error) {
	if f, ok := loggerFactories[name]; ok {
		return f(conf)
	}
	return nil, fmt.Errorf("logger: unknown driver %q", name)
}
