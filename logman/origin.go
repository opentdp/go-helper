package logman

import (
	"log/slog"
	"os"

	"github.com/opentdp/go-helper/onquit"
)

var (
	Debug = slog.Debug
	Info  = slog.Info
	Warn  = slog.Warn
	Error = slog.Error
)

func Fatal(msg string, args ...any) {

	onquit.CallQuitFuncs() // 调用所有退出函数

	Error(msg, args...)
	os.Exit(1)

}
