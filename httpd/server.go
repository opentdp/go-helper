package httpd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentdp/go-helper/logman"
)

func Server(addr string) {

	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	// 以协程方式启用监听，防止阻塞后续的中断信号处理
	go func() {
		logman.Info("server starting", "address", addr)
		if err := server.ListenAndServe(); err != nil {
			logman.Warn("server terminated", "error", err)
		}
		logman.Info("server exited")
		os.Exit(0)
	}()

	// 创建监听中断信号通道
	quit := make(chan os.Signal, 1)
	// SIGTERM: `kill`
	// SIGINT : `kill -2` 或 CTRL + C
	// SIGKILL: `kill -9`，无法捕获，故而不做处理
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// 等待信号，如果没有则保持阻塞
	<-quit

	logman.Warn("server exiting...")

	// 创建一个剩余5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器而不中断任何活动连接
	if err := server.Shutdown(ctx); err != nil {
		logman.Fatal("server forced to shutdown", "error", err)
	}

}
