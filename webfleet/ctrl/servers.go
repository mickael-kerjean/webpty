package ctrl

import (
	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/webfleet/model"
	"net/http"
)

type serverInstance struct {
	id       string
	isOnline bool
	info     map[string]string
}

func ListServers(res http.ResponseWriter, req *http.Request) {
	Log.Info("SERVER %+v", model.Machines.List())

	// list, err := Webfleet.List() // []ServerInstance,error
	// if err != nil {
	// 	ErrorPage(res, err, 400)
	// 	return
	// }
	// RenderTemplate("tmpl/dashbard.html", list)
}
