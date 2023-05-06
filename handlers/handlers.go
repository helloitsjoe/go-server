package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Name     string `json:"name" binding:"required,max=1000"`
	Password string `json:"password" binding:"required"`
	// Age      int    `json:"age" binding:"required,gte=1,lte=150"`
}

type Handlers struct {
	users map[uuid.UUID]User
}

func NewHandlers() *Handlers {
	users := make(map[uuid.UUID]User)
	return &Handlers{users}
}

func (h Handlers) Home(c *gin.Context) {
	c.File("index.html")
}

func (h Handlers) Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h Handlers) GetUser(c *gin.Context) {
	id := c.Param("id")
	foo := c.GetHeader("X-Foo")
	c.Header("X-Foo-Response", foo)
	c.String(http.StatusOK, "Hello %s %s", id, foo)
}

func (h Handlers) GetAllUsers(c *gin.Context) {
	usersArr := []User{}
	for _, v := range h.users {
		usersArr = append(usersArr, v)
	}
	jsonUsers, err := json.Marshal(usersArr)
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.JSON(http.StatusOK, string(jsonUsers))
}

func (h Handlers) Register(c *gin.Context) {
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	h.users[id] = User{body.Name, body.Password}
	c.JSON(http.StatusOK, gin.H{"name": body.Name, "password": body.Password, "id": id})
}
