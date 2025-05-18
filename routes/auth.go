package routes

import (
	"github.com/gin-gonic/gin"
	"temp/controllers"
	"temp/middleware"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.GET("/myprofile", middleware.AuthMiddleware(), controllers.MyProfile)
		auth.PUT("/myprofile", middleware.AuthMiddleware(), controllers.UpdateMyProfile)
		
		auth.GET("/match", middleware.AuthMiddleware(), controllers.MatchUsers)
		auth.POST("/send-request", middleware.AuthMiddleware(), controllers.SendSkillRequest)
		auth.POST("/respond-request", middleware.AuthMiddleware(), controllers.RespondToSkillRequest)
	}
}





