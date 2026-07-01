package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/starfork/stargo/logger"
)

var LevelNames = map[slog.Level]string{
	slog.LevelDebug: "debug",
	slog.LevelInfo:  "info",
	slog.LevelWarn:  "warn",
	slog.LevelError: "error",
}

func init() {
	logger.Register("slog", NewSlogLogger)
}

type SlogLogger struct {
	l *slog.Logger
}

func NewSlogLogger(conf *logger.Config) (logger.Logger, error) {
	level := slog.LevelInfo
	switch conf.Level {
	case -1:
		level = slog.LevelDebug
	case 0:
		level = slog.LevelInfo
	case 1:
		level = slog.LevelWarn
	case 2:
		level = slog.LevelError
	}

	var handler slog.Handler
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	if conf.Target == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	}

	return &SlogLogger{l: slog.New(handler)}, nil
}

func (s *SlogLogger) log(level slog.Level, format string, v ...any) {
	if len(v) > 0 {
		s.l.Log(context.Background(), level, fmt.Sprintf(format, v...))
	} else {
		s.l.Log(context.Background(), level, format)
	}
}

func (s *SlogLogger) Warnf(format string, v ...any)  { s.log(slog.LevelWarn, format, v...) }
func (s *SlogLogger) Debugf(format string, v ...any) { s.log(slog.LevelDebug, format, v...) }
func (s *SlogLogger) Errorf(format string, v ...any) { s.log(slog.LevelError, format, v...) }
func (s *SlogLogger) Fatalf(format string, v ...any) { s.log(slog.LevelError, format, v...); os.Exit(1) }
func (s *SlogLogger) Infof(format string, v ...any)  { s.log(slog.LevelInfo, format, v...) }
func (s *SlogLogger) String() string                  { return "slog" }
func (s *SlogLogger) Options() logger.Options         { return logger.Options{} }

func (s *SlogLogger) WithContext(ctx context.Context) logger.Logger {
	return s
}
