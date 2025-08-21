package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"taskmanager/database"
	"taskmanager/models"
)

type CreateTaskTypeInput struct {
	Label string `json:"label" binding:"required"`
}

type UpdateTaskTypeInput struct {
	Label string `json:"label"`
}

// CreateTaskType creates a new task type
func CreateTaskType(c *gin.Context) {
	var input CreateTaskTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskType := models.TaskType{Label: input.Label}

	if err := database.DB.Create(&taskType).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"errors": []string{"Task type with this label already exists"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to create task type"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": taskType})
}

// GetTaskTypes retrieves all task types
func GetTaskTypes(c *gin.Context) {
	var taskTypes []models.TaskType
	if err := database.DB.Find(&taskTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": taskTypes})
}

// GetTaskTypeByID retrieves a single task type by ID
func GetTaskTypeByID(c *gin.Context) {
	id := c.Param("id")
	var taskType models.TaskType

	if err := database.DB.Where("id = ?", id).First(&taskType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task type not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": taskType})
}

// UpdateTaskType updates an existing task type
func UpdateTaskType(c *gin.Context) {
	id := c.Param("id")
	var taskType models.TaskType

	if err := database.DB.Where("id = ?", id).First(&taskType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task type not found"})
		return
	}

	var input UpdateTaskTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&taskType).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": taskType})
}

// DeleteTaskType deletes a task type
func DeleteTaskType(c *gin.Context) {
	id := c.Param("id")
	var taskType models.TaskType

	if err := database.DB.Where("id = ?", id).First(&taskType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task type not found"})
		return
	}

	database.DB.Delete(&taskType)

	c.JSON(http.StatusOK, gin.H{"message": "Task type deleted successfully"})
}
