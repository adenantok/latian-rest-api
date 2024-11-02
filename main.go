package main

import (
	"latian-rest-api/config"
	"latian-rest-api/controllers"
	"latian-rest-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()        // Connect to the database
	models.Migrate(config.DB) // Run migrations using the DB instance

	router := gin.Default()

	bukuService := controllers.NewBukuService()                  // Create instance of BukuService
	bukuController := controllers.NewBukuController(bukuService) // Create instance of BukuController

	router.GET("/buku", bukuController.GetBukuHandler) // Use handler for GET request

	log.Fatal(router.Run(":3000")) // Start the server
}