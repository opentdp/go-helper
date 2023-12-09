package httpd

import (
	"context"
	"net/http"
	"time"

	"github.com/opentdp/go-helper/logman"
	"github.com/opentdp/go-helper/onquit"
)

var server *http.Server

func Server(addr string, options ...any) {

	if engine == nil {
		if len(options) > 0 {
			Engine(options[0].(bool))
		} else {
			Engine(false)
		}
	}

	server = &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	onquit.Register(func() {
		// 创建一个剩余15秒超时的上下文
		logman.Warn("httpd will close within 9 seconds")
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()

		// 优雅地关闭服务而不中断任何活动连接
		if err := server.Shutdown(ctx); err != nil {
			logman.Warn("httpd forced to close", "error", err)
			server.Close()
		}
	})

	logman.Info("httpd start", "address", addr)
	if err := server.ListenAndServe(); err != nil {
		logman.Warn(err.Error())
	}

}
