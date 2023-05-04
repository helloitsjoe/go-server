package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// middleware/auth
// cors

var allowedHeaders = []string{
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

func middy(c *gin.Context) {
	fmt.Println("Middleware!")
	c.Next()
}

// func myCors(c *gin.Context) {
// 	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
// 	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
// 	c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
// 	c.Header("Access-Control-Allow-Credentials", "true")

// 	if c.Request.Method == "OPTIONS" {
// 		c.AbortWithStatus(204)
// 		return
// 	}

// 	c.Next()
// }

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
	router.Use(middy)
	// router.Use(myCors)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", home)
	router.GET("/ping", pong)
	router.GET("/user/:id", getUser)
	router.POST("/data", postData)
	fmt.Println("Listening on http://localhost:8080")
	router.Run()
}
