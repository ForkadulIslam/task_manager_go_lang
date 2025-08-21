package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanager/database"
	"taskmanager/models"
)

type AssignTaskToGroupInput struct {
	GroupID uint `json:"group_id" binding:"required"`
	TaskID  uint `json:"task_id" binding:"required"`
}

// AssignTaskToGroup assigns a task to a group
func AssignTaskToGroup(c *gin.Context) {
	var input AssignTaskToGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assign := models.AssignTaskToGroup{GroupID: input.GroupID, TaskID: input.TaskID}

	if err := database.DB.Create(&assign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": assign})
}

// RemoveTaskAssignmentFromGroup removes a task assignment from a group
func RemoveTaskAssignmentFromGroup(c *gin.Context) {
	id := c.Param("id")
	var assign models.AssignTaskToGroup

	if err := database.DB.Where("id = ?", id).First(&assign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task assignment not found"})
		return
	}

	database.DB.Delete(&assign)

	c.JSON(http.StatusOK, gin.H{"message": "Task assignment removed successfully"})
}

// GetTasksAssignedToGroup retrieves all tasks assigned to a specific group
func GetTasksAssignedToGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	var assigns []models.AssignTaskToGroup

	if err := database.DB.Where("group_id = ?", groupID).Find(&assigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks assigned to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assigns})
}

// GetGroupsAssignedToTask retrieves all groups a specific task is assigned to
func GetGroupsAssignedToTask(c *gin.Context) {
	taskID := c.Param("task_id")
	var assigns []models.AssignTaskToGroup

	if err := database.DB.Where("task_id = ?", taskID).Find(&assigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups assigned to task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assigns})
}
