package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)

var AllowedMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
var AllowedHeaders = []string{
	"accept",
	"origin",
	"X-CSRF-Token",
	"Authorization",
	"Cache-Control",
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-Requested-With",
}

func AuthMiddlware(c *gin.Context) {
	// token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)

	token, err := c.Cookie("token")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if token != "fake-auth-token" {
		fmt.Println("Not authorized!")
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("token", token)

	fmt.Println("Authorized!")
	c.Next()
}

func MyCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Methods", strings.Join(AllowedMethods, ","))
	c.Header("Access-Control-Allow-Headers", strings.Join(AllowedHeaders, ","))
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
