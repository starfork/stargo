package logger

import (
	"os"

	"time"

	"github.com/starfork/stargo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjackv2 "gopkg.in/natefinch/lumberjack.v2"
)

const (
	Console = "console"
	File    = "file"
)

var (
	Leavel = zap.DebugLevel
	Target = Console
)

// NewEncoderConfig new
func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// TimeEncoder time format
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Init Init
func New(conf *config.LogConfig) *zap.SugaredLogger {

	var writeSyncer zapcore.WriteSyncer
	target := conf.Target
	if target == "" {
		target = Console
	}

	if target == Console {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	if target == File {
		w := zapcore.AddSync(&lumberjackv2.Logger{
			Filename:   conf.LogFile,
			MaxSize:    1024, // megabytes
			MaxBackups: 10,
			MaxAge:     7, // days
		})

		writeSyncer = zapcore.NewMultiWriteSyncer(w)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writeSyncer,
		Leavel,
	)

	return zap.New(core, zap.AddCaller()).Sugar()
}
