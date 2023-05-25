package httpd

import (
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
