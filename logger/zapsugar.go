package logger

import (
	"os"
	"time"

	"github.com/starfork/stargo/config"
	"github.com/starfork/stargo/util/ustring"
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
	level := zap.DebugLevel

	conf := &config.LogConfig{
		Target:  "console",
		LogFile: "debug.log",
	}
	if len(c) > 0 && c[0] != nil {
		tmp := c[0]
		conf.Target = ustring.OrString("console", tmp.Target)
		conf.LogFile = ustring.OrString("debug.log", tmp.LogFile)
		if tmp.Level > 0 {
			level = zapcore.Level(conf.Level)
		}
		//
	}
	if conf.Target == "console" {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	if conf.Target == "file" {
		w := zapcore.AddSync(&lumberjackv2.Logger{
			Filename:   conf.LogFile,
			MaxSize:    1024, // megabytes
			MaxBackups: 10,
			MaxAge:     7, // days
		})

		writeSyncer = zapcore.NewMultiWriteSyncer(w)
	}
	//fmt.Println(level)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writeSyncer,
		//zap.DebugLevel,
		level,
	)
	return zap.New(core, zap.AddCaller()).Sugar()
}
