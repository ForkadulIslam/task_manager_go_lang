package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanager/database"
	"taskmanager/models"
)

func GetNotifications(c *gin.Context) {
	userID := uint(c.MustGet("user_id").(float64))
	var notifications []models.Notification
	database.DB.Preload("User").Preload("Task").Where("user_id = ?", userID).Order("created_at desc").Find(&notifications)
	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}
	notification.IsRead = true
	database.DB.Save(&notification)
	c.JSON(http.StatusOK, notification)
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID := uint(c.MustGet("user_id").(float64))
	database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}
