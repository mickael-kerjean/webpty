package middleware

import "net/http"

func WithAuth(fn http.HandlerFunc) http.HandlerFunc {
	return fn
}
