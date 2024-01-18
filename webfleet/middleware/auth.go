package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	. "github.com/mickael-kerjean/webpty/common"
	"github.com/mickael-kerjean/webpty/webfleet/view"
)

var (
	AUTH_DRIVER string
	AUTH_FUNC   func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
)

func init() {
	AUTH_DRIVER = strings.ToLower(os.Getenv("AUTH_DRIVER"))
	switch AUTH_DRIVER {
	case "yolo":
		AUTH_FUNC = driverYolo
	case "simple":
		AUTH_FUNC = driverSimple
	default:
		AUTH_FUNC = nil
		Log.Error("missing AUTH_DRIVER env variable")
		os.Exit(1)
	}
}

func WithAuth(fn http.HandlerFunc) http.HandlerFunc {
	if AUTH_FUNC == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			view.ErrorPage(w, errors.New("Missing authenticator"), http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		AUTH_FUNC(w, r, fn)
	}
}
