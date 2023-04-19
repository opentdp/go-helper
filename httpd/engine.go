package httpd

import (
	"mime"

	"github.com/gin-gonic/gin"
	"github.com/open-tdp/go-helper/logman"
)

func Engine(debug bool) *gin.Engine {

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = logman.AutoWriter("gin-access")
	gin.DefaultErrorWriter = logman.AutoWriter("gin-error")

	// 初始化
	return gin.Default()

}

func init() {

	// 重写文件类型
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")

}
