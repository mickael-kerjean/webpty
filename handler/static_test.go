package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleStatic(t *testing.T) {
	for _, test := range []struct {
		Method       string
		Url          string
		ExpectedCode int
	}{
		{"GET", "/", 200},
		{"GET", "/wp-admin/", 404},
		{"GET", "/node_modules/xterm/lib/xterm.js", 200},
		{"GET", "/node_modules/xterm/LICENSE", 200},
		{"GET", "/node_modules/xterm/does_not_exist", 404},
		{"GET", "/app/app.css", 200},
	} {
		req, _ := http.NewRequest(test.Method, test.Url, nil)
		res := httptest.NewRecorder()

		HandleStatic(res, req)
		assert.Equal(t, test.ExpectedCode, res.Code)
	}
}
