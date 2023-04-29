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
	router.GET("/user/:id", controllers.GetUser)
	router.GET("/users", controllers.GetUsers)
	router.DELETE("/user/:id", controllers.DeleteUser)
	router.PATCH("/user/:id", controllers.UpdateUser)
	return router
}