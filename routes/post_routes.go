package routes

import (
	"bn/controllers"
	"bn/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	posts := router.Group("/api/posts")
	posts.Use(middleware.AuthMiddleware()) // Apply JWT authentication to these routes
	{
		posts.POST("/", controllers.CreatePost())
		posts.GET("/", controllers.GetPosts())
		posts.GET("/:id", controllers.GetPostByID())
		posts.PUT("/:id", controllers.UpdatePost())
		posts.DELETE("/:id", controllers.DeletePost())
	}
}