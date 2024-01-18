package middleware

import (
	"net/http"
)

func driverYolo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
