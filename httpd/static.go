package httpd

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func Static(urlPrefix, root string) {

	hfs := gin.Dir(root, false)
	engine.Use(StaticServe(urlPrefix, hfs))

}

func StaticIndex(urlPrefix, root string) {

	hfs := gin.Dir(root, true)
	engine.Use(StaticServe(urlPrefix, hfs))

}

func StaticEmbed(urlPrefix, sub string, efs *embed.FS) {

	var hfs http.FileSystem

	if sub == "" {
		hfs = http.FS(efs)
	} else {
		eb, _ := fs.Sub(efs, sub)
		hfs = http.FS(eb)
	}

	engine.Use(StaticServe(urlPrefix, hfs))

}

func StaticServe(urlPrefix string, hfs http.FileSystem) gin.HandlerFunc {

	fileserver := http.FileServer(hfs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	isExists := func(prefix, filepath string) bool {
		if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
			if f, err := hfs.Open(path.Join("/", p)); err == nil {
				defer f.Close()
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		if isExists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}

}
