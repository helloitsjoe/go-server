package main

import (
	"encoding/json"
	"fmt"
	"go_server/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var inMemoryUsers = map[uuid.UUID]User{}

func home(c *gin.Context) {
	c.File("index.html")
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	foo := c.GetHeader("X-Foo")
	c.Header("X-Foo-Response", foo)
	c.String(http.StatusOK, "Hello %s %s", id, foo)
}

func getAllUsers(c *gin.Context) {
	usersArr := []User{}
	for _, v := range inMemoryUsers {
		usersArr = append(usersArr, v)
	}
	jsonUsers, err := json.Marshal(usersArr)
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.JSON(http.StatusOK, string(jsonUsers))
}

type User struct {
	Name     string `json:"name" binding:"required,max=1000"`
	Password string `json:"password" binding:"required"`
	// Age      int    `json:"age" binding:"required,gte=1,lte=150"`
}

func register(c *gin.Context) {
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	inMemoryUsers[id] = User{body.Name, body.Password}
	c.JSON(http.StatusOK, gin.H{"name": body.Name, "password": body.Password, "id": id})
}

func main() {
	router := gin.Default()
	// router.Use(myCors)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     middleware.AllowedMethods,
		AllowHeaders:     middleware.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", home)
	router.GET("/ping", pong)
	router.GET("/user/:id", middleware.AuthMiddlware, getUser)
	router.GET("/users", getAllUsers)
	// router.POST("/data", authMiddlware, postData)
	router.POST("/register", register)
	// router.POST("/login", login)
	fmt.Println("Listening on http://localhost:8080")
	router.Run()
}
