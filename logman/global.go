package logman

import (
	"os"

	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level    string
	Target   string
	Filename string
}

func SetDefault(args *Config) {

	var level slog.Level
	var handler slog.Handler

	level.UnmarshalText([]byte(args.Level))

	opt := slog.HandlerOptions{
		Level: level,
	}

	switch args.Target {
	case "file":
		fw := fileWriter(args.Filename)
		handler = opt.NewJSONHandler(fw)
	case "stdout":
		handler = opt.NewTextHandler(os.Stdout)
	default:
		handler = opt.NewTextHandler(os.Stderr)
	}

	slog.SetDefault(slog.New(handler))

}

func fileWriter(filename string) *lumberjack.Logger {

	return &lumberjack.Logger{
		Compress:   true,     // 是否压缩/归档旧文件
		Filename:   filename, // 日志文件位置
		MaxSize:    100,      // 单个日志文件最大值(单位：MB)
		MaxBackups: 21,       // 保留旧文件的最大个数
		MaxAge:     7,        // 保留旧文件的最大天数
	}

}
