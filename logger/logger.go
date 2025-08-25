package logger

type Logger interface {
	// Logf writes a formatted log entry

	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
	Infof(format string, v ...any)

	// String returns the name of logger
	String() string

	Options() Options
}

var DefaultLogger Logger = NewLogger()
