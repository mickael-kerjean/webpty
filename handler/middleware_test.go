package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	for _, test := range []struct {
		AuthHeader   string
		ExpectedCode int
	}{
		{"Basic dXNlcm5hbWU6cGFzc3dvcmQ=", http.StatusUnauthorized},
		{"", http.StatusUnauthorized},
	} {
		req, _ := http.NewRequest("GET", "/", nil)
		if test.AuthHeader != "" {
			req.Header.Set("Authorization", test.AuthHeader)
		}
		res := httptest.NewRecorder()
		BasicAuth(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte("OK"))
		})(res, req)
		assert.Equal(t, test.ExpectedCode, res.Code)
	}
}
