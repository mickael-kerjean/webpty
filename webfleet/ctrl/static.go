package ctrl

import (
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/view"
	"io"
	"net/http"
)

func ServeFile(path string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		file, err := view.Tmpl.Open(path)
		if err != nil {
			ctrl.ErrorPage(res, ErrNotFound, 404)
			return
		}
		io.Copy(res, file)
	}
}
