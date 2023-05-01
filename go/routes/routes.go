package routes

import (
	"jwt/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB) (*gin.Engine){
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
	router.GET("/welcome", controllers.Welcome)
	return router
}