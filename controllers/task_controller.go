package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"taskmanager/database"
	"taskmanager/models"
	"taskmanager/utils"
)

type CreateTaskInput struct {
	Label       string         `json:"label" binding:"required,min=3,max=255"`
	TaskTypeID  uint           `json:"task_type_id" binding:"required,gt=0"`
	Priority    string         `json:"priority" binding:"required,oneof=Normal Medium High Escalation"`
	StartDate   utils.Date `json:"start_date" binding:"required"`
	DueDate     *utils.Date `json:"due_date" binding:"omitempty,gtefield=StartDate"`
	Description string         `json:"description"`
	Attachment  string         `json:"attachment"`
	Status      string         `json:"status"`
}

type UpdateTaskInput struct {
	Label       string         `json:"label" binding:"min=3,max=255"`
	TaskTypeID  uint           `json:"task_type_id" binding:"gt=0"`
	Priority    string         `json:"priority" binding:"oneof=Normal Medium High Escalation"`
	StartDate   utils.Date `json:"start_date" binding:"required"`
	DueDate     *utils.Date `json:"due_date" binding:"omitempty,gtefield=StartDate"`
	Description string         `json:"description"`
	Attachment  string         `json:"attachment"`
	Status      string         `json:"status"`
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var input CreateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			_ = ve.Translate(trans)

			for _, e := range ve {
				switch e.Field() {
				case "Label":
					if e.Tag() == "required" {
						errors = append(errors, "Task label is required")
					} else if e.Tag() == "min" {
						errors = append(errors, "Task label must be at least 3 characters long")
					} else if e.Tag() == "max" {
						errors = append(errors, "Task label cannot exceed 255 characters")
					}
				case "TaskTypeID":
					if e.Tag() == "required" {
						errors = append(errors, "Task type is required")
					} else if e.Tag() == "gt" {
						errors = append(errors, "Task type ID must be greater than 0")
					}
				case "Priority":
					if e.Tag() == "required" {
						errors = append(errors, "Priority is required")
					} else if e.Tag() == "oneof" {
						errors = append(errors, "Invalid priority value. Must be Normal, Medium, High, or Escalation")
					}
				case "StartDate":
					if e.Tag() == "required" {
						errors = append(errors, "Start date is required")
					}
				case "DueDate":
					if e.Tag() == "gtefield" {
						errors = append(errors, "Due date must be greater than or equal to start date")
					}
				default:
					errors = append(errors, e.Translate(trans))
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// Check if the task type exists
	var taskType models.TaskType
	if err := database.DB.First(&taskType, input.TaskTypeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid task type ID"}})
		return
	}

	// Placeholder CreatedBy - should come from authenticated user
	task := models.Task{
		Label:       input.Label,
		TaskTypeID:  input.TaskTypeID,
		Priority:    input.Priority,
		StartDate:   input.StartDate.Time,
		Description: input.Description,
		Attachment:  input.Attachment,
		Status:      input.Status,
		CreatedBy:   uint(c.MustGet("user_id").(float64)),
	}

	if input.DueDate != nil {
		task.DueDate = &input.DueDate.Time
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

// GetTasks retrieves all tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := database.DB.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTaskByID retrieves a single task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			_ = ve.Translate(trans)

			for _, e := range ve {
				switch e.Field() {
				case "Label":
					if e.Tag() == "min" {
						errors = append(errors, "Task label must be at least 3 characters long")
					} else if e.Tag() == "max" {
						errors = append(errors, "Task label cannot exceed 255 characters")
					}
				case "TaskTypeID":
					if e.Tag() == "gt" {
						errors = append(errors, "Task type ID must be greater than 0")
					}
				case "Priority":
					if e.Tag() == "oneof" {
						errors = append(errors, "Invalid priority value. Must be Normal, Medium, High, or Escalation")
					}
				case "DueDate":
					if e.Tag() == "gtefield" {
						errors = append(errors, "Due date must be greater than or equal to start date")
					}
				default:
					errors = append(errors, e.Translate(trans))
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	database.DB.Model(&task).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// DeleteTask deletes a task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	database.DB.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}