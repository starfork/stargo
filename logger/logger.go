package logger

type Logger interface {
	// Logf writes a formatted log entry

	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})

	// String returns the name of logger
	String() string

	Options() Options
}

var DefaultLogger Logger = NewLogger()
