package logger

import (
	"context"
	"fmt"
	"maps"
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

func (l *defaultLogger) Fields(fields map[string]any) Logger {
	l.Lock()
	nfields := make(map[string]any, len(l.opts.Fields))

	maps.Copy(nfields, l.opts.Fields)
	l.Unlock()

	maps.Copy(nfields, fields)

	return &defaultLogger{opts: Options{
		Level: l.opts.Level,
	}}
}

func copyFields(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	maps.Copy(dst, src)

	return dst
}

func (l *defaultLogger) Logf(level Level, format string, v ...any) {
	fmt.Printf(format, v...)
}

func (l *defaultLogger) Warnf(format string, v ...any) {
	l.Logf(WarnLevel, format, v...)
}
func (l *defaultLogger) Debugf(format string, v ...any) {
	l.Logf(DebugLevel, format, v...)
}
func (l *defaultLogger) Errorf(format string, v ...any) {
	l.Logf(ErrorLevel, format, v...)
}
func (l *defaultLogger) Fatalf(format string, v ...any) {
	l.Logf(FatalLevel, format, v...)
	os.Exit(1)
}
func (l *defaultLogger) Infof(format string, v ...any) {
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
		Fields:          make(map[string]any),
		Out:             os.Stderr,
		CallerSkipCount: 2,
		Context:         context.Background(),
	}

	l := &defaultLogger{opts: options}
	if err := l.Init(opts...); err != nil {
		l.Logf(FatalLevel, "init logger fail %s", err.Error())
	}

	return l
}
