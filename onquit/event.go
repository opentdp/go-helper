package onquit

import (
	"os"
	"os/signal"
	"syscall"
)

var quit chan os.Signal
var onQuitFuncs []func()

// 注册退出时的回调函数
func Register(onExit func()) {

	onQuitFuncs = append(onQuitFuncs, onExit)

	// 避免重复注册断信号通道
	if quit != nil {
		return
	}

	// 创建监听中断信号通道
	quit = make(chan os.Signal, 1)

	// SIGTERM: `kill`
	// SIGINT : `kill -2` 或 CTRL + C
	// SIGKILL: `kill -9`，无法捕获，故而不做处理
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// 等待退出信号
	go func() {
		<-quit
		CallQuitFuncs()
		os.Exit(0) // 退出
	}()

}

// 调用所有退出函数
func CallQuitFuncs() {

	for _, fn := range onQuitFuncs {
		fn()
	}

}
