package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// query params
// sub routes
// headers
// middleware/auth
// cors
// status codes
// post body

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
	c.String(http.StatusOK, "Hello %s", id)
}

type BasicPostBody struct {
	// Fields must be exported (capitalized)
	Name string
	Age  int
}

func postData(c *gin.Context) {
	var body BasicPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
