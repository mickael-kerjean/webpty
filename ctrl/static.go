package ctrl

import (
	"embed"
	. "github.com/mickael-kerjean/webpty/common"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

//go:generate bash static.sh
//go:embed src
var efs embed.FS

func HandleStatic(res http.ResponseWriter, req *http.Request) {
	urlPath := req.URL.Path
	if urlPath == "/" {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		ServeFile(res, req, "src/index.html")
		return
	}

	actualPath := ""
	if strings.HasPrefix(urlPath, "/app/") {
		base := filepath.Base(urlPath)
		actualPath = "src/app/" + base
	} else if strings.HasPrefix(urlPath, "/node_modules/") {
		actualPath = "src" + urlPath
	} else {
		ErrorPage(res, ErrNotFound, 404)
		return
	}
	ServeFile(res, req, actualPath)

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

// gzip compression was done using:
// find . -type f | grep -v '\.ts$' | grep -v '\.map$' | xargs gzip -f -k --best
func ServeFile(res http.ResponseWriter, req *http.Request, filePath string) {
	head := res.Header()
	acceptEncoding := req.Header.Get("Accept-Encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		if file, err := efs.Open(filePath + ".gz"); err == nil {
			head.Set("Content-Encoding", "gzip")
			io.Copy(res, file)
			file.Close()
			return
		}
	}

	file, err := efs.Open(filePath)
	if err != nil {
		ErrorPage(res, ErrNotFound, 404)
		return
	}

	if strings.HasSuffix(filePath, ".js") {
		res.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(filePath, ".css") {
		res.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasSuffix(filePath, ".html") {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	io.Copy(res, file)
	file.Close()
}

//go:embed src/favicon.ico
var IconFavicon []byte

func ServeFavicon(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write(IconFavicon)
}

func HealthCheck(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
	res.WriteHeader(200)
}
