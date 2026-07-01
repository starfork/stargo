package zap

import (
	"context"

	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/util/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Register("zap", NewZapLogger)
}

type ZapLogger struct {
	l *zap.SugaredLogger
}

func NewZapLogger(conf *logger.Config) (logger.Logger, error) {
	level := zapcore.InfoLevel
	switch conf.Level {
	case -1:
		level = zapcore.DebugLevel
	case 0:
		level = zapcore.InfoLevel
	case 1:
		level = zapcore.WarnLevel
	case 2:
		level = zapcore.ErrorLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	if conf.Target == "console" {
		cfg.Encoding = "console"
	}

	l, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{l: l.Sugar()}, nil
}

func (z *ZapLogger) Warnf(format string, v ...any)  { z.l.Warnf(format, v...) }
func (z *ZapLogger) Debugf(format string, v ...any) { z.l.Debugf(format, v...) }
func (z *ZapLogger) Errorf(format string, v ...any) { z.l.Errorf(format, v...) }
func (z *ZapLogger) Fatalf(format string, v ...any) { z.l.Fatalf(format, v...) }
func (z *ZapLogger) Infof(format string, v ...any)  { z.l.Infof(format, v...) }
func (z *ZapLogger) String() string                  { return "zap" }
func (z *ZapLogger) Options() logger.Options         { return logger.Options{} }

func (z *ZapLogger) WithContext(ctx context.Context) logger.Logger {
	traceID := tracing.TraceIDFromContext(ctx)
	spanID := tracing.SpanIDFromContext(ctx)
	return &ZapLogger{l: z.l.With(
		"trace_id", traceID,
		"span_id", spanID,
	)}
}
