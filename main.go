package main

import (
	"time" // New import for time.Duration
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // New import

	"taskmanager/database"
	"taskmanager/routes"
)

func main() {
	database.ConnectDatabase()

	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend origin
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
