package logger

var DefaultLogger Logger = NewZapSugar()

type Logger interface {
	Debugf(template string, args ...interface{})
}
