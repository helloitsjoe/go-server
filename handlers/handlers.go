package handlers

import (
	"fmt"
	"net/http"

	"go_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string    `json:"name" binding:"required,max=1000"`
	Password string    `json:"password" binding:"required"`
	Id       uuid.UUID `json:"id"`
	// Age      int    `json:"age" binding:"required,gte=1,lte=150"`
}

type UserMap map[uuid.UUID]User

type Handlers struct {
	users UserMap
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashed), err
}

func (h *Handlers) SeedUsers(users UserMap) {
	if users != nil {
		for _, user := range users {
			password, err := hashPassword(user.Password)
			if err != nil {
				panic(err)
			}
			user.Password = password
			h.users[user.Id] = user
		}
	} else {
		h.users = UserMap{}
	}
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

type UserResponse struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

func (h Handlers) GetAllUsers(c *gin.Context) {
	usersArr := []UserResponse{}
	for _, v := range h.users {
		usersArr = append(usersArr, UserResponse{v.Name, v.Id})
	}
	c.JSON(http.StatusOK, usersArr)
}

func (h Handlers) Register(c *gin.Context) {
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	hashedPass, err := hashPassword(body.Password)
	if err != nil {
		c.AbortWithStatus(500)
	}

	fmt.Println(hashedPass)
	h.users[id] = User{body.Name, string(hashedPass), id}
	fmt.Println("register", h.users)
	c.JSON(http.StatusOK, gin.H{"name": body.Name, "id": id})
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

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.Password))
	if err != nil {
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
