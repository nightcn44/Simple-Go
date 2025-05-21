package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"

	"backend/config"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	r := gin.Default()
	routes.UserRoute(r)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
