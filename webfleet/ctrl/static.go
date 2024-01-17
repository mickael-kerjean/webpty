package ctrl

import (
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/view"
	"io"
	"net/http"
)

func ServeFile(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := view.Tmpl.Open(path)
		if err != nil {
			ctrl.ErrorPage(w, ErrNotFound, 404)
			return
		}
		io.Copy(w, file)
	}
}
