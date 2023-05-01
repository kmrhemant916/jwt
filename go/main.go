package main

import (
	"jwt/models"
	"jwt/routes"
)


func main() {
	db,err := models.Connection()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Token{})
	router := routes.SetupRoutes(db)
	router.Run()
}