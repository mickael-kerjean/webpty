package view

import (
	"embed"
)

//go:embed *.html
var Tmpl embed.FS
