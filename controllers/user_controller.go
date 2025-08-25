package controllers

import (
	"net/http"
	"taskmanager/database"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
)

// GetUsers retrieves all active users with the 'User' label
func GetUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}
