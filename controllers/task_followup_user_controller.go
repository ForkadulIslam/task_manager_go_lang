package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanager/database"
	"taskmanager/models"
)

type AddTaskFollowupUserInput struct {
	UserID  uint `json:"user_id" binding:"required"`
	TaskID  uint `json:"task_id" binding:"required"`
	Remarks string `json:"remarks"`
}

// AddTaskFollowupUser adds a user for task follow-up
func AddTaskFollowupUser(c *gin.Context) {
	var input AddTaskFollowupUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followupUser := models.TaskFollowupUser{UserID: input.UserID, TaskID: input.TaskID, Remarks: input.Remarks}

	if err := database.DB.Create(&followupUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task follow-up user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": followupUser})
}

// RemoveTaskFollowupUser removes a task follow-up user
func RemoveTaskFollowupUser(c *gin.Context) {
	id := c.Param("id")
	var followupUser models.TaskFollowupUser

	if err := database.DB.Where("id = ?", id).First(&followupUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task follow-up user not found"})
		return
	}

	database.DB.Delete(&followupUser)

	c.JSON(http.StatusOK, gin.H{"message": "Task follow-up user removed successfully"})
}

// GetFollowupUsersForTask retrieves all follow-up users for a specific task
func GetFollowupUsersForTask(c *gin.Context) {
	taskID := c.Param("task_id")
	var followupUsers []models.TaskFollowupUser

	if err := database.DB.Where("task_id = ?", taskID).Find(&followupUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve follow-up users for task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": followupUsers})
}

// GetTasksFollowedByUser retrieves all tasks a specific user is following up on
func GetTasksFollowedByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var followupUsers []models.TaskFollowupUser

	if err := database.DB.Where("user_id = ?", userID).Find(&followupUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks followed by user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": followupUsers})
}
