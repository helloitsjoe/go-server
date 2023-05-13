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

type User struct {
	Name     string
	Password string
	Id       uuid.UUID
}

func TestRegister(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"Name": "foo", "Password": "bar"}`))
	w := makeRequest(http.MethodPost, "/register", reader, nil)

	var dat User
	err := json.Unmarshal(w.Body.Bytes(), &dat)
	if err != nil {
		t.Log("error", err)
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dat.Name, "foo")
	assert.Equal(t, dat.Password, "bar")
	assert.IsType(t, "string", dat.Id.String())
}

func TestUsers(t *testing.T) {
	id := uuid.New()
	seedUsers := handlers.UserMap{}
	seedUsers[id] = handlers.User{Name: "foo", Password: "bar", Id: id}
	w := makeRequest(http.MethodGet, "/users", nil, seedUsers)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf(`[{"name":"foo","password":"bar","id":"%s"}]`, id), w.Body.String())
}
