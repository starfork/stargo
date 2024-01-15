package zapsugar

import (
	"os"
	"time"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/util/ustring"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjackv2 "gopkg.in/natefinch/lumberjack.v2"
)

type ZapSugar struct {
	opts  logger.Options
	sugar *zap.SugaredLogger
}

// NewEncoderConfig new
func encoderConfig() zapcore.EncoderConfig {
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
func NewZapSugar(c ...*logger.Config) logger.Logger {
	var writeSyncer zapcore.WriteSyncer
	level := zap.DebugLevel

	conf := &logger.Config{
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
		zapcore.NewConsoleEncoder(encoderConfig()),
		writeSyncer,
		//zap.DebugLevel,
		level,
	)
	return &ZapSugar{
		sugar: zap.New(core, zap.AddCaller()).Sugar(),
	}
}

func (e *ZapSugar) Log(level logger.Level, v ...interface{}) {
	e.sugar.Info(v...)
}

// Logf writes a formatted log entry
func (e *ZapSugar) Logf(level logger.Level, format string, v ...interface{}) {
	e.sugar.Infof(format, v...)
}

func (e *ZapSugar) Debugf(format string, v ...interface{}) {
	e.sugar.Debugf(format, v...)
}
func (e *ZapSugar) Fatalf(format string, v ...interface{}) {
	e.sugar.Debugf(format, v...)
	os.Exit(1)
}
func (e *ZapSugar) Infof(format string, v ...interface{}) {}

// String returns the name of logger
func (e *ZapSugar) String() string {
	return "zapsugar"
}

func (e *ZapSugar) Options() logger.Options {
	return e.opts
}
