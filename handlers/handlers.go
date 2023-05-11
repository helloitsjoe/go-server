package handlers

import (
	"encoding/json"
	"net/http"

	"go_server/utils"

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
	c.File("static/index.html")
}

func (h Handlers) Pong(c *gin.Context) {
	foo := c.GetHeader("X-Foo")
	c.Header("X-Foo-Response", foo)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h Handlers) GetUser(c *gin.Context) {
	if c.GetString("token") == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.File("static/user.html")
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

func (h Handlers) LoginGet(c *gin.Context) {
	c.File("static/login.html")
}

func (h Handlers) LoginPost(c *gin.Context) {
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var foundUser User
	for _, user := range h.users {
		if user.Name == body.Name {
			foundUser = user
		}
	}

	if foundUser == (User{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if foundUser.Password != body.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	token, err := utils.GenerateToken(foundUser.Name)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.SetCookie("token", token, 1000*60*60, "/", "http://localhost:8080", true, true)

	c.JSON(http.StatusOK, gin.H{"name": foundUser.Name})
}
