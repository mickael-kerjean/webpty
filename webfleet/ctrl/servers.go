package ctrl

import (
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/ctrl"
	"github.com/mickael-kerjean/webpty/webfleet/model"
	"github.com/mickael-kerjean/webpty/webfleet/view"
	"html/template"
	"net/http"
)

func ListServers(w http.ResponseWriter, r *http.Request) {
	lsrv, err := model.Machines.List()
	if err != nil {
		Log.Error("ctrl::servers list machines failed '%s'", err.Error())
		ctrl.ErrorPage(w, err, 400)
		return
	}
	tmpl, err := template.ParseFS(view.Tmpl, "dashboard.html")
	if err != nil {
		ctrl.ErrorPage(w, err, 500)
		return
	}
	tmpl.Execute(
		w,
		map[string]interface{}{
			"servers": lsrv,
		},
	)
}
