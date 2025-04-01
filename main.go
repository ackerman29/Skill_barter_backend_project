package main

import (
	"github.com/gin-gonic/gin"
	"temp/config"
	"temp/middleware"
	"temp/routes"
	"github.com/gin-contrib/cors"
	// "temp/helpers"
)

func main() {
	config.ConnectDB()
	r := gin.Default()
	r.Use(cors.Default())
	// Public route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "PONG"})
	})

	// 🔒 Protected route
	r.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		email := c.MustGet("email").(string)
		c.JSON(200, gin.H{
			"message": "Welcome to protected route!",
			"email":   email,
		})
	})
	routes.AuthRoutes(r)
	
	r.Run(":8000")

}
