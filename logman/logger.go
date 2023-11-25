package logman

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	name   string
	logger *slog.Logger
	ctx    context.Context
}

func Named(name string) *Logger {

	return &Logger{
		name: name, logger: NewLogger(name),
	}

}

func (l *Logger) log(level slog.Level, msg string, args ...any) {

	args = append([]any{"Logger", l.name}, args...)
	l.logger.Log(l.ctx, level, msg, args...)

}

func (l *Logger) Debug(msg string, args ...any) {

	l.log(slog.LevelDebug, msg, args...)

}

func (l *Logger) Info(msg string, args ...any) {

	l.log(slog.LevelInfo, msg, args...)

}

func (l *Logger) Warn(msg string, args ...any) {

	l.log(slog.LevelWarn, msg, args...)

}

func (l *Logger) Error(msg string, args ...any) {

	l.log(slog.LevelError, msg, args...)

}

func (l *Logger) Fatal(msg string, args ...any) {

	l.log(slog.LevelError, msg, args...)
	os.Exit(1)

}
