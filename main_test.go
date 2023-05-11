package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	r.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	w := setup(http.MethodGet, "/ping", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestHome(t *testing.T) {
	w := setup(http.MethodGet, "/", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "<title>Go Server</title>"))
}
