package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/webfleet/view"
)

var simple_users = map[string]string{}

func init() {
	if AUTH_DRIVER != "simple" {
		return
	}

	users := strings.Split(os.Getenv("AUTH_USER"), ",")
	for _, userpass := range users {
		v := strings.SplitN(userpass, ":", 2)
		if len(v) != 2 {
			continue
		}
		simple_users[v[0]] = v[1]
	}
	if len(simple_users) == 0 {
		Log.Error("You don't have any users setup: eg: AUTH_USER=test:test")
		os.Exit(1)
		return
	}
}

func driverSimple(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username, password, ok := r.BasicAuth()
	if !ok || (simple_users[username] == "" || simple_users[username] != password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		view.ErrorPage(w, errors.New("Not Authorised"), http.StatusUnauthorized)
		return
	}
	next(w, r)
}
