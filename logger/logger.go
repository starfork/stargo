package logger

type Logger interface {
	Log(level Level, v ...interface{})
	// Logf writes a formatted log entry
	Logf(level Level, format string, v ...interface{})

	Debugf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})

	// String returns the name of logger
	String() string

	Options() Options
}

var DefaultLogger Logger = NewLogger()
