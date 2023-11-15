package onquit

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/opentdp/go-helper/logman"
)

var quit chan os.Signal
var onQuitFuncs []func()

func Register(onExit func()) {

	onQuitFuncs = append(onQuitFuncs, onExit)

	// 避免重复注册断信号通道
	if quit != nil {
		return
	}

	// 创建监听中断信号通道
	quit := make(chan os.Signal, 1)

	// SIGTERM: `kill`
	// SIGINT : `kill -2` 或 CTRL + C
	// SIGKILL: `kill -9`，无法捕获，故而不做处理
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// 接收到退出信号时，遍历 onExitFuncs 切片，并调用每个函数
	go func() {
		<-quit
		logman.Warn("exiting...")
		for _, fn := range onQuitFuncs {
			fn()
		}
		os.Exit(0)
	}()

}
