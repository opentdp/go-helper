package httpd

import (
	"context"
	"net/http"
	"time"

	"github.com/opentdp/go-helper/logman"
	"github.com/opentdp/go-helper/onquit"
)

func Server(addr string) {

	if engine == nil {
		Engine(false)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	onquit.Register(func() {
		// 创建一个剩余15秒超时的上下文
		logman.Warn("httpd will be closed, wait 15 seconds")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		// 优雅地关闭服务器而不中断任何活动连接
		if err := server.Shutdown(ctx); err != nil {
			logman.Warn("httpd forced to close", "error", err)
			server.Close()
		}
	})

	logman.Info("httpd starting", "address", addr)
	if err := server.ListenAndServe(); err != nil {
		logman.Warn(err.Error())
	}

}
