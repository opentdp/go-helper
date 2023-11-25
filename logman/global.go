package logman

import (
	"io"
	"log/slog"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

var config = &Config{}

type Config struct {
	Level    string `note:"日志级别 DEBUG|INFO|WARN|ERROR"`
	Target   string `note:"日志输出设备 both|file|stdout|stderr"`
	Storage  string `note:"日志文件存储目录"`
	Filename string `note:"默认日志文件名"`
}

func SetDefault(args *Config) {

	config.Level = args.Level
	config.Target = args.Target
	config.Storage = args.Storage

	slog.SetDefault(NewLogger(args.Filename))

}

func NewLogger(name string) *slog.Logger {

	var handler slog.Handler
	var level slog.Level = 0

	level.UnmarshalText([]byte(config.Level))

	option := &slog.HandlerOptions{
		Level: level,
	}

	writer := AutoWriter(name)
	handler = slog.NewTextHandler(writer, option)

	return slog.New(handler)

}

func AutoWriter(name string) io.Writer {

	switch config.Target {
	case "file":
		return FileWriter(name)
	case "both":
		return io.MultiWriter(os.Stdout, FileWriter(name))
	case "stderr":
		return os.Stderr
	default:
		return os.Stdout
	}

}

func FileWriter(name string) *lumberjack.Logger {

	f := path.Join(config.Storage, name) + ".log"

	return &lumberjack.Logger{
		Compress:   true, // 是否压缩/归档旧文件
		Filename:   f,    // 日志文件位置
		MaxSize:    100,  // 单个日志文件最大值(单位：MB)
		MaxBackups: 21,   // 保留旧文件的最大个数
		MaxAge:     7,    // 保留旧文件的最大天数
	}

}
