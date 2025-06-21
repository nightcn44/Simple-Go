package main

import (
	"log"
	"os"

	"bn/config"
	"bn/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	config.ConnectDB()
	defer config.DB.Disconnect(nil) // Disconnect when main function exits

	router := gin.Default()

	// Setup routes
	routes.AuthRoutes(router)
	routes.PostRoutes(router) // Protected routes for posts

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}
