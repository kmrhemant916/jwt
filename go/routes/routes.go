package routes

import (
	"jwt/controllers"
	"jwt/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	HEADER = "x-auth-token"
)

func SetupRoutes(db *gorm.DB) (*gin.Engine){
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middlewares.AuthHeader(HEADER))
	protectedRoutes.POST("/logout", controllers.Logout)
	protectedRoutes.GET("/welcome", controllers.Welcome)
	return router
}