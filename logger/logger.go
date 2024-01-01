package logger

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
}
