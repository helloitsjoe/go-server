package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_server/handlers"
	"go_server/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func makeRequest(req *http.Request, seedUsers handlers.UserMap) *httptest.ResponseRecorder {
	r := setupRouter(seedUsers, true)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := makeRequest(req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestHome(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := makeRequest(req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "<title>Go Server</title>"))
}

func TestRegister(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "bar"}`))
	req, _ := http.NewRequest(http.MethodPost, "/register", reader)
	w := makeRequest(req, nil)

	var dat handlers.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &dat)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dat.Name, "foo")
	assert.IsType(t, "string", dat.Id.String())
}

func TestUsers(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := makeRequest(req, seedUsers)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf(`[{"name":"foo","id":"%s"}]`, id), w.Body.String())
}

func TestUser(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	token, _ := utils.GenerateToken(seedUsers[id].Name)
	req, _ := http.NewRequest(http.MethodGet, "/user/foo", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	w := makeRequest(req, nil)

	t.Log(w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), "<body>"))
}

// ===== Login Tests =====

func TestLoginPostSuccess(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "bar"}`))

	req, _ := http.NewRequest(http.MethodPost, "/login", reader)
	w := makeRequest(req, seedUsers)
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

	req, _ := http.NewRequest(http.MethodPost, "/login", reader)
	w := makeRequest(req, seedUsers)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginPostNoPassword(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "foo"}`))
	req, _ := http.NewRequest(http.MethodPost, "/login", reader)
	w := makeRequest(req, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostNoName(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Password": "foo"}`))
	req, _ := http.NewRequest(http.MethodPost, "/login", reader)
	w := makeRequest(req, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostEmpty(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/login", nil)
	w := makeRequest(req, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginPostNoUser(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "bar", "Password": "foo"}`))
	req, _ := http.NewRequest(http.MethodPost, "/login", reader)
	w := makeRequest(req, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
