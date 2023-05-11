package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRequest(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	r.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	w := makeRequest(http.MethodGet, "/ping", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestHome(t *testing.T) {
	w := makeRequest(http.MethodGet, "/", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "<title>Go Server</title>"))
}

func TestRegister(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"name": "foo", "password": "bar"}`))
	w := makeRequest(http.MethodPost, "/register", reader)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.IsType(t, "string", w.Body.String())
}

func TestUsers(t *testing.T) {
	w := makeRequest(http.MethodGet, "/users", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[]`, w.Body.String())
}
