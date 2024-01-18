package view

import (
	"errors"
	"net/http"
	"text/template"
)

func ErrorPage(res http.ResponseWriter, err error, code int) {
	if err == nil {
		err = errors.New("Unknown error")
	}
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
