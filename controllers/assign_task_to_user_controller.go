package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanager/database"
	"taskmanager/models"
)

type AssignTaskToUserInput struct {
	UserID uint `json:"user_id" binding:"required"`
	TaskID uint `json:"task_id" binding:"required"`
}

// AssignTaskToUser assigns a task to a user
func AssignTaskToUser(c *gin.Context) {
	var input AssignTaskToUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// Check if UserID exists
	var user models.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"User not found"}})
		return
	}

	// Check if TaskID exists
	var task models.Task
	if err := database.DB.First(&task, input.TaskID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Task not found"}})
		return
	}

	// Check for duplicate assignment
	var existingAssignment models.AssignTaskToUser
	if err := database.DB.Where("user_id = ? AND task_id = ?", input.UserID, input.TaskID).First(&existingAssignment).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"errors": []string{"Task already assigned to this user"}})
		return
	}

	assign := models.AssignTaskToUser{UserID: input.UserID, TaskID: input.TaskID}

	if err := database.DB.Create(&assign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to assign task to user"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": assign})
}

// RemoveTaskAssignmentFromUser removes a task assignment from a user
func RemoveTaskAssignmentFromUser(c *gin.Context) {
	id := c.Param("id")
	var assign models.AssignTaskToUser

	if err := database.DB.Where("id = ?", id).First(&assign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task assignment not found"})
		return
	}

	database.DB.Delete(&assign)

	c.JSON(http.StatusOK, gin.H{"message": "Task assignment removed successfully"})
}

// GetTasksAssignedToUser retrieves all tasks assigned to a specific user
func GetTasksAssignedToUser(c *gin.Context) {
	userID := c.Param("user_id")
	var assigns []models.AssignTaskToUser

	if err := database.DB.Where("user_id = ?", userID).Find(&assigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks assigned to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assigns})
}

// GetUsersAssignedToTask retrieves all users assigned to a specific task
func GetUsersAssignedToTask(c *gin.Context) {
	taskID := c.Param("task_id")
	var assigns []models.AssignTaskToUser

	if err := database.DB.Where("task_id = ?", taskID).Find(&assigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users assigned to task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assigns})
}
