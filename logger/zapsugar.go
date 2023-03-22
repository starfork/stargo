package logger

import (
	"os"
	"time"

	"github.com/starfork/stargo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjackv2 "gopkg.in/natefinch/lumberjack.v2"
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
func NewZapSugar(c ...*config.LogConfig) *zap.SugaredLogger {

	var writeSyncer zapcore.WriteSyncer
	var target, logFile string

	if len(c) > 0 {
		target = c[0].Target
		logFile = c[0].LogFile
	}
	if target == "" {
		target = "console"
	}
	if logFile == "" {
		logFile = "debug.log"
	}

	if target == "console" {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	if target == "file" {
		w := zapcore.AddSync(&lumberjackv2.Logger{
			Filename:   logFile,
			MaxSize:    1024, // megabytes
			MaxBackups: 10,
			MaxAge:     7, // days
		})

		writeSyncer = zapcore.NewMultiWriteSyncer(w)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writeSyncer,
		zap.DebugLevel,
	)

	return zap.New(core, zap.AddCaller()).Sugar()
}
