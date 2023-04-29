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

func home() {
	// TODO: Serve HTML
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

func main() {
	router := gin.Default()
	router.GET("/ping", pong)
	router.GET("/user/:id", getUser)
	fmt.Println("Listening on http://localhost:8080")
	router.Run()
}
