package handlers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSeedUsersEmpty(t *testing.T) {
	handlers := NewHandlers()
	handlers.SeedUsers(nil)
	assert.Equal(t, handlers.users, UserMap{})
}

func TestSeedUsers(t *testing.T) {
	id := uuid.New()
	seed := UserMap{}
	seed[id] = User{Name: "foo", Password: "bar", Id: id}
	handlers := NewHandlers()
	handlers.SeedUsers(seed)

	expected := UserMap{}
	expected[id] = User{Name: "foo", Password: "bar", Id: id}

	assert.Equal(t, len(expected), len(handlers.users))
	assert.Equal(t, expected[id].Id, handlers.users[id].Id)
	assert.Equal(t, expected[id].Name, handlers.users[id].Name)
	assert.NotEqual(t, expected[id].Password, handlers.users[id].Password)
}
