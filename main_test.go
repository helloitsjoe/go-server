package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_server/handlers"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func makeRequest(method string, url string, body io.Reader, seedUsers handlers.UserMap) *httptest.ResponseRecorder {
	r := setupRouter(seedUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	r.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	w := makeRequest(http.MethodGet, "/ping", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestHome(t *testing.T) {
	w := makeRequest(http.MethodGet, "/", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "<title>Go Server</title>"))
}

func TestRegister(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "bar"}`))
	w := makeRequest(http.MethodPost, "/register", reader, nil)

	var dat handlers.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &dat)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dat.Name, "foo")
	assert.IsType(t, "string", dat.Id.String())
}

func TestLoginPostSuccess(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "bar"}`))

	w := makeRequest(http.MethodPost, "/login", reader, seedUsers)
	cookie := w.Result().Cookies()[0]

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"name":"foo"}`, w.Body.String())
	assert.IsType(t, "string", cookie.Value)
	assert.Equal(t, true, cookie.HttpOnly)
	assert.Equal(t, true, cookie.Secure)
	// assert.Equal(t, "http://localhost:8080", cookie.Domain)
}

func TestLoginPostUnauthorized(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "not-bar"}`))

	w := makeRequest(http.MethodPost, "/login", reader, seedUsers)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginPostNoPassword(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "foo"}`))
	w := makeRequest(http.MethodPost, "/login", reader, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostNoName(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Password": "foo"}`))
	w := makeRequest(http.MethodPost, "/login", reader, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostEmpty(t *testing.T) {
	w := makeRequest(http.MethodPost, "/login", nil, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostNoUser(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "bar", "Password": "foo"}`))
	w := makeRequest(http.MethodPost, "/login", reader, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUsers(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	w := makeRequest(http.MethodGet, "/users", nil, seedUsers)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf(`[{"name":"foo","id":"%s"}]`, id), w.Body.String())
}
