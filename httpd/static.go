package httpd

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func Static(prefix, root string) {

	hfs := gin.Dir(root, false)
	engine.Use(StaticServe(prefix, hfs))

}

func StaticIndex(prefix, root string) {

	hfs := gin.Dir(root, true)
	engine.Use(StaticServe(prefix, hfs))

}

func StaticEmbed(prefix, sub string, efs *embed.FS) {

	var hfs http.FileSystem

	if sub == "" {
		hfs = http.FS(efs)
	} else {
		sub, _ := fs.Sub(efs, sub)
		hfs = http.FS(sub)
	}

	engine.Use(StaticServe(prefix, hfs))

}

func StaticServe(prefix string, hfs http.FileSystem) gin.HandlerFunc {

	fileServer := http.FileServer(hfs)
	if prefix != "" {
		fileServer = http.StripPrefix(prefix, fileServer)
	}

	isExists := func(p, s string) bool {
		if p := strings.TrimPrefix(s, p); len(p) < len(s) {
			if f, err := hfs.Open(path.Join("/", p)); err == nil {
				defer f.Close()
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		if isExists(prefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}

}
