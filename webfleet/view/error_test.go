package view

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorPageWithError(t *testing.T) {
	w := httptest.NewRecorder()
	ErrorPage(w, errors.New("__error__"), http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestErrorPageWithoutError(t *testing.T) {
	w := httptest.NewRecorder()
	ErrorPage(w, nil, http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, w.Code)
}
