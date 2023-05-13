package main

import (
	"fmt"
	"go_server/handlers"
	"go_server/middleware"
	"go_server/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var PORT = utils.GetEnv("PORT", "8080")

func setupRouter(usersSeed handlers.UserMap) *gin.Engine {
	h := handlers.NewHandlers()

	h.SeedUsers(usersSeed)

	router := gin.Default()
	router.SetTrustedProxies(nil)
	// router.Use(middleware.MyCors)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     middleware.AllowedMethods,
		AllowHeaders:     middleware.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Static("static", "./static")

	router.GET("/", h.Home)
	router.GET("/ping", h.Pong)
	router.GET("/user/:id", middleware.AuthMiddlware, h.GetUser)
	router.GET("/users", h.GetAllUsers)
	// router.POST("/data", authMiddlware, postData)
	router.POST("/register", h.Register)
	router.GET("/login", h.LoginGet)
	router.POST("/login", h.LoginPost)

	return router
}

func main() {
	r := setupRouter(nil)
	fmt.Printf("Listening on http://localhost:%s\n", PORT)
	r.Run(fmt.Sprintf(":%s", PORT))
}
