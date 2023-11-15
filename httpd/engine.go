package httpd

import (
	"mime"

	"github.com/gin-gonic/gin"

	"github.com/opentdp/go-helper/logman"
)

var engine *gin.Engine

func Engine(debug bool) *gin.Engine {

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = logman.AutoWriter("gin-access")
	gin.DefaultErrorWriter = logman.AutoWriter("gin-error")

	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")

	engine = gin.Default()

	return engine

}

func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {

	return engine.Group(relativePath, handlers...)

}

func Use(middleware ...gin.HandlerFunc) gin.IRoutes {

	return engine.Use(middleware...)

}
