package handler

import (
	"embed"
	. "github.com/mickael-kerjean/virtualshell/common"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed src
var efs embed.FS

func HandleStatic(res http.ResponseWriter, req *http.Request) {
	urlPath := req.URL.Path
	if urlPath == "/" {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		f, err := efs.Open("src/index.html")
		if err != nil {
			Log.Error("handler::static efs.Open %s", err.Error())
			ErrorPage(res, ErrNotFound, 404)
			return
		}
		io.Copy(res, f)
		return
	}

	var (
		f   fs.File
		err error
	)
	if strings.HasPrefix(urlPath, "/app/") {
		base := filepath.Base(urlPath)
		f, err = efs.Open("src/app/" + base)
	} else if strings.HasPrefix(urlPath, "/node_modules/") {
		f, err = efs.Open("src" + urlPath)
	} else {
		ErrorPage(res, ErrNotFound, 404)
		return
	}

	if err != nil {
		Log.Error("handler::static fs.Open url[%s] err[%s]", urlPath, err.Error())
		ErrorPage(res, ErrNotFound, 404)
		return
	}

	if strings.HasSuffix(urlPath, ".js") {
		res.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(urlPath, ".css") {
		res.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	io.Copy(res, f)
}

func ErrorPage(res http.ResponseWriter, err error, code int) {
	res.WriteHeader(code)
	tmpl, _ := template.
		New("handler::static").
		Parse(`<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>{{ .ErrorCode }}</title>
  </head>
  <body>
    <div>
        <h1>{{ .ErrorCode }}</h1>
        <p> {{ .ErrorMessage }}</p>
    </div>
    <style>
      body { font-family: monaco, monospace; background: #010088; color: #f8ffff;
             display: flex; align-items: center; justify-content: center;
             text-align: center; }
      body > div { margin-top: -50px; }
      h1 { font-size: 150px; line-height: 150px; margin: 0; background: #f8ffff;
           color: #010088; padding: 5px 20px; }
      p { font-size: 25px; line-height: 25px; }
    </style>
  </body>
</html>`)
	tmpl.Execute(res, struct {
		ErrorCode    int
		ErrorMessage string
	}{code, err.Error()})
}
