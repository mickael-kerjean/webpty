package ctrl

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

func TestEndpointHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	HealthCheck(w, httptest.NewRequest("GET", "/healthz", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestEndpointFavicon(t *testing.T) {
	w := httptest.NewRecorder()
	ServeFavicon(w, httptest.NewRequest("GET", "/favicon.ico", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, len(w.Body.String()) > 10)
}
