package main

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplicationBootOK(t *testing.T) {
	// given
	go main()
	time.Sleep(100 * time.Millisecond)

	// when
	resp, err := (&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}).Get("https://127.0.0.1:3456/healthz")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// then
	srv.Close()
}
