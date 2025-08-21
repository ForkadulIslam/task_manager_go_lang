package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanager/database"
	"taskmanager/models"
)

type AddUserToGroupInput struct {
	UserID  uint `json:"user_id" binding:"required"`
	GroupID uint `json:"group_id" binding:"required"`
}

// AddUserToGroup adds a user to a group
func AddUserToGroup(c *gin.Context) {
	var input AddUserToGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userGroup := models.UserGroup{UserID: input.UserID, GroupID: input.GroupID}

	if err := database.DB.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": userGroup})
}

// RemoveUserFromGroup removes a user from a group
func RemoveUserFromGroup(c *gin.Context) {
	id := c.Param("id")
	var userGroup models.UserGroup

	if err := database.DB.Where("id = ?", id).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User-group association not found"})
		return
	}

	database.DB.Delete(&userGroup)

	c.JSON(http.StatusOK, gin.H{"message": "User removed from group successfully"})
}

// GetUsersInGroup retrieves all users in a specific group
func GetUsersInGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	var userGroups []models.UserGroup

	if err := database.DB.Where("group_id = ?", groupID).Find(&userGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users in group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userGroups})
}

// GetGroupsForUser retrieves all groups a specific user belongs to
func GetGroupsForUser(c *gin.Context) {
	userID := c.Param("user_id")
	var userGroups []models.UserGroup

	if err := database.DB.Where("user_id = ?", userID).Find(&userGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userGroups})
}
