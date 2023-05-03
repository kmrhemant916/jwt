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
	models.DatabaseMigration(db)
	router := routes.SetupRoutes(db)
	router.Run()
}