package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// middleware/auth
// cors

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

type BasicPostBody struct {
	Name string `json:"name" binding:"required,max=1000"`
	Age  int    `json:"age" binding:"required,gte=1,lte=150"`
}

func postData(c *gin.Context) {
	var body BasicPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}

func main() {
	router := gin.Default()
	router.GET("/", home)
	router.GET("/ping", pong)
	router.GET("/user/:id", getUser)
	router.POST("/data", postData)
	fmt.Println("Listening on http://localhost:8080")
	router.Run()
}
