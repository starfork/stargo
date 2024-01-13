package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
)

func init() {
	lvl, err := GetLevel(os.Getenv("STARGO_LOG_LEVEL"))
	if err != nil {
		lvl = InfoLevel
	}

	DefaultLogger = NewLogger(WithLevel(lvl))
}

type defaultLogger struct {
	opts Options
	sync.RWMutex
}

// Init (opts...) should only overwrite provided options.
func (l *defaultLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}

	return nil
}

func (l *defaultLogger) String() string {
	return "default"
}

func (l *defaultLogger) Fields(fields map[string]interface{}) Logger {
	l.Lock()
	nfields := make(map[string]interface{}, len(l.opts.Fields))

	for k, v := range l.opts.Fields {
		nfields[k] = v
	}
	l.Unlock()

	for k, v := range fields {
		nfields[k] = v
	}

	return &defaultLogger{opts: Options{
		Level: l.opts.Level,
	}}
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func (l *defaultLogger) Log(level Level, v ...interface{}) {
	fmt.Printf("  %v\n", v)
}

func (l *defaultLogger) Logf(level Level, format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.Logf(DebugLevel, format, v...)
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.Logf(InfoLevel, format, v...)
}

func (l *defaultLogger) Options() Options {
	// not guard against options Context values
	l.RLock()
	defer l.RUnlock()

	opts := l.opts
	opts.Fields = copyFields(l.opts.Fields)

	return opts
}

// NewLogger builds a new logger based on options.
func NewLogger(opts ...Option) Logger {
	// Default options
	options := Options{
		Level:           InfoLevel,
		Fields:          make(map[string]interface{}),
		Out:             os.Stderr,
		CallerSkipCount: 2,
		Context:         context.Background(),
	}

	l := &defaultLogger{opts: options}
	if err := l.Init(opts...); err != nil {
		l.Log(FatalLevel, err)
	}

	return l
}
