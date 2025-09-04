package main

import (
	"time" // New import for time.Duration

	"github.com/gin-contrib/cors" // New import
	"github.com/gin-gonic/gin"

	"taskmanager/database"
	"taskmanager/routes"
)

func main() {
	database.ConnectDatabase()

	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static files from the "uploads" directory
	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r)
	r.Run(":8080")
}
