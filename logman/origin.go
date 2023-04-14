package logman

import (
	"os"

	"golang.org/x/exp/slog"
)

var (
	Debug = slog.Debug
	Info  = slog.Info
	Warn  = slog.Warn
	Error = slog.Error
)

func Fatal(msg string, args ...any) {

	Error(msg, args...)
	os.Exit(1)

}
