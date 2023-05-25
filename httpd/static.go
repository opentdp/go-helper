package httpd

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
)

func Static(urlPrefix, dir string) {

	lf := static.LocalFile(dir, false)
	engine.Use(static.Serve(urlPrefix, lf))

}

func StaticIndex(urlPrefix, dir string) {

	lf := static.LocalFile(dir, true)
	engine.Use(static.Serve(urlPrefix, lf))

}

func StaticEmbed(urlPrefix, dir string, efs *embed.FS) {

	ui, _ := fs.Sub(efs, dir)
	engine.StaticFS(urlPrefix, http.FS(ui))

}
