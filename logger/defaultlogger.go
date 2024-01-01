package logger

import "fmt"

type DefaultLog struct{}

func (e *DefaultLog) Debugf(template string, args ...interface{}) {
	fmt.Printf(template, args...)

}
func (e *DefaultLog) Infof(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

var DefaultLogger = &DefaultLog{}
