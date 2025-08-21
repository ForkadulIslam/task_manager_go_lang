package main

import (
	"github.com/gin-gonic/gin"

	"taskmanager/database"
	"taskmanager/routes"
)

func main() {
	database.ConnectDatabase()

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
