package ctrl

import (
	"embed"
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/ctrl"
	"io"
	"net/http"
)

func ServeFile(fs embed.FS, path string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		file, err := fs.Open(path)
		if err != nil {
			ctrl.ErrorPage(res, ErrNotFound, 404)
			return
		}
		io.Copy(res, file)
	}
}
