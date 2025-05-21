package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
)

func UserRoute(r *gin.Engine) {
	r.GET("/users", controllers.GetUsers)
	r.POST("/user", controllers.CreateUser)
}
