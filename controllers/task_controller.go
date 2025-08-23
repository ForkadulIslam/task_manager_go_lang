package controllers

import (
	"fmt"
	"net/http"

	"taskmanager/database"
	"taskmanager/models"
	"taskmanager/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateTaskInput struct {
	Label            string      `json:"label" binding:"required,min=3,max=255"`
	TaskTypeID       uint        `json:"task_type_id" binding:"required,gt=0"`
	Priority         string      `json:"priority" binding:"required,oneof=Normal Medium High Escalation"`
	StartDate        utils.Date  `json:"start_date" binding:"required"`
	DueDate          *utils.Date `json:"due_date" binding:"omitempty,gtefield=StartDate"`
	Description      string      `json:"description"`
	Attachment       string      `json:"attachment"`
	Status           string      `json:"status"`
	AssignedToUsers  []uint      `json:"assigned_to_users" binding:"omitempty,dive,gt=0"`
	AssignedToGroups []uint      `json:"assigned_to_groups" binding:"omitempty,dive,gt=0"`
	FollowUpUsers    []uint      `json:"follow_up_users" binding:"omitempty,dive,gt=0"`
}

type UpdateTaskInput struct {
	Label            string      `json:"label" binding:"min=3,max=255"`
	TaskTypeID       uint        `json:"task_type_id" binding:"gt=0"`
	Priority         string      `json:"priority" binding:"oneof=Normal Medium High Escalation"`
	StartDate        utils.Date  `json:"start_date" binding:"required"`
	DueDate          *utils.Date `json:"due_date" binding:"omitempty,gtefield=StartDate"`
	Description      string      `json:"description"`
	Attachment       string      `json:"attachment"`
	Status           string      `json:"status"`
	AssignedToUsers  []uint      `json:"assigned_to_users" binding:"omitempty,dive,gt=0"`
	AssignedToGroups []uint      `json:"assigned_to_groups" binding:"omitempty,dive,gt=0"`
	FollowUpUsers    []uint      `json:"follow_up_users" binding:"omitempty,dive,gt=0"`
}

type UpdateTaskStatusInput struct {
	Status string `json:"status" binding:"required,oneof=Pending 'In Progress' 'In Review' Completed"`
}

type AddTaskCommentInput struct {
	Comment string `json:"comment" binding:"required"`
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

	// Validate AssignedToUsers
	for _, userID := range input.AssignedToUsers {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("User with ID %d not found", userID)}})
			return
		}
	}

	// Validate AssignedToGroups
	for _, groupID := range input.AssignedToGroups {
		var group models.Group
		if err := database.DB.First(&group, groupID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Group with ID %d not found", groupID)}})
			return
		}
	}

	// Validate FollowUpUsers
	for _, userID := range input.FollowUpUsers {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Follow-up user with ID %d not found", userID)}})
			return
		}
	}

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

	// Assign task to users
	for _, userID := range input.AssignedToUsers {
		assignToUser := models.AssignTaskToUser{
			TaskID: task.ID,
			UserID: userID,
		}
		if err := database.DB.Create(&assignToUser).Error; err != nil {
			// Handle error, perhaps rollback task creation or log it
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to user"})
			return
		}
	}

	// Assign task to groups
	for _, groupID := range input.AssignedToGroups {
		assignToGroup := models.AssignTaskToGroup{
			TaskID:  task.ID,
			GroupID: groupID,
		}
		if err := database.DB.Create(&assignToGroup).Error; err != nil {
			// Handle error, perhaps rollback task creation or log it
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to group"})
			return
		}
	}

	// Assign task to follow-up users
	for _, userID := range input.FollowUpUsers {
		followUpUser := models.TaskFollowupUser{
			TaskID: task.ID,
			UserID: userID,
		}
		if err := database.DB.Create(&followUpUser).Error; err != nil {
			// Handle error, perhaps rollback task creation or log it
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to follow-up user"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

// GetTasks retrieves all tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := database.DB.Preload("AssignedUsers.User").Preload("AssignedGroups.Group.Users").Preload("FollowupUsers.User").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTaskByID retrieves a single task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.Preload("AssignedUsers.User").Preload("AssignedGroups.Group.Users").Preload("FollowupUsers.User").Preload("Comments.User").Where("id = ?", id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// isUserAssigned checks if a user is assigned to a task, either directly or through a group
func isUserAssigned(db *gorm.DB, userID uint, taskID uint) (bool, error) {
	// Check for direct assignment
	var directAssignment int64
	err := db.Model(&models.AssignTaskToUser{}).Where("task_id = ? AND user_id = ?", taskID, userID).Count(&directAssignment).Error
	if err != nil {
		return false, err
	}
	if directAssignment > 0 {
		return true, nil
	}

	// Check for group assignment
	var groupAssignment int64
	err = db.Model(&models.UserGroup{}).
		Joins("JOIN assign_task_to_groups ON user_groups.group_id = assign_task_to_groups.group_id").
		Where("assign_task_to_groups.task_id = ? AND user_groups.user_id = ?", taskID, userID).
		Count(&groupAssignment).Error
	if err != nil {
		return false, err
	}

	return groupAssignment > 0, nil
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	authUserID := uint(c.MustGet("user_id").(float64))
	if task.CreatedBy != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You are not authorized to update this task"}})
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
				// You can customize error messages here as in CreateTask
				errors = append(errors, e.Translate(trans))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// Validate TaskTypeID, if provided
	if input.TaskTypeID != 0 {
		var taskType models.TaskType
		if err := database.DB.First(&taskType, input.TaskTypeID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid task type ID"}})
			return
		}
	}

	// Validate AssignedToUsers, if provided
	if input.AssignedToUsers != nil {
		for _, userID := range input.AssignedToUsers {
			var user models.User
			if err := database.DB.First(&user, userID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("User with ID %d not found", userID)}})
				return
			}
		}
	}

	// Validate AssignedToGroups, if provided
	if input.AssignedToGroups != nil {
		for _, groupID := range input.AssignedToGroups {
			var group models.Group
			if err := database.DB.First(&group, groupID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Group with ID %d not found", groupID)}})
				return
			}
		}
	}

	// Validate FollowUpUsers, if provided
	if input.FollowUpUsers != nil {
		for _, userID := range input.FollowUpUsers {
			var user models.User
			if err := database.DB.First(&user, userID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Follow-up user with ID %d not found", userID)}})
				return
			}
		}
	}

	// Using a transaction to ensure atomicity
	tx := database.DB.Begin()

	// Update task fields
	if err := tx.Model(&task).Updates(models.Task{
		Label:       input.Label,
		TaskTypeID:  input.TaskTypeID,
		Priority:    input.Priority,
		StartDate:   input.StartDate.Time,
		DueDate:     &input.DueDate.Time,
		Description: input.Description,
		Attachment:  input.Attachment,
		Status:      input.Status,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update associations
	if input.AssignedToUsers != nil {
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.AssignTaskToUser{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assigned users"})
			return
		}
		for _, userID := range input.AssignedToUsers {
			assignToUser := models.AssignTaskToUser{TaskID: task.ID, UserID: userID}
			if err := tx.Create(&assignToUser).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to user"})
				return
			}
		}
	}

	if input.AssignedToGroups != nil {
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.AssignTaskToGroup{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assigned groups"})
			return
		}
		for _, groupID := range input.AssignedToGroups {
			assignToGroup := models.AssignTaskToGroup{TaskID: task.ID, GroupID: groupID}
			if err := tx.Create(&assignToGroup).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to group"})
				return
			}
		}
	}

	if input.FollowUpUsers != nil {
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskFollowupUser{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update follow-up users"})
			return
		}
		for _, userID := range input.FollowUpUsers {
			followUpUser := models.TaskFollowupUser{TaskID: task.ID, UserID: userID}
			if err := tx.Create(&followUpUser).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to follow-up user"})
				return
			}
		}
	}


	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// UpdateTaskStatus updates the status of a task by an assigned user
func UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	authUserID := uint(c.MustGet("user_id").(float64))

	// Authorization check
	assigned, err := isUserAssigned(database.DB, authUserID, task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user assignment"})
		return
	}
	if !assigned {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You are not authorized to update the status of this task"}})
		return
	}

	var input UpdateTaskStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")
			_ = ve.Translate(trans)
			for _, e := range ve {
				errors = append(errors, e.Translate(trans))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// Using a transaction to ensure atomicity
	tx := database.DB.Begin()

	// Update task status
	if err := tx.Model(&task).Update("status", input.Status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	// Log the status update
	statusLog := models.TaskStatusUpdateLog{
		TaskID: task.ID,
		UserID: authUserID,
		Status: input.Status,
	}
	if err := tx.Create(&statusLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log status update"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully"})
}

// isUserAssignedOrFollowup checks if a user is assigned to a task (directly or via group) or is a followup user.
func isUserAssignedOrFollowup(db *gorm.DB, userID uint, taskID uint) (bool, error) {
	// Check if assigned (directly or via group)
	assigned, err := isUserAssigned(db, userID, taskID)
	if err != nil {
		return false, err
	}
	if assigned {
		return true, nil
	}

	// Check for followup user assignment
	var followupAssignment int64
	err = db.Model(&models.TaskFollowupUser{}).Where("task_id = ? AND user_id = ?", taskID, userID).Count(&followupAssignment).Error
	if err != nil {
		return false, err
	}

	return followupAssignment > 0, nil
}

// AddTaskComment adds a comment to a task
func AddTaskComment(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	authUserID := uint(c.MustGet("user_id").(float64))

	// Authorization check
	allowed, err := isUserAssignedOrFollowup(database.DB, authUserID, task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user assignment"})
		return
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You are not authorized to comment on this task"}})
		return
	}

	var input AddTaskCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	comment := models.TaskCommentLog{
		TaskID:  task.ID,
		UserID:  authUserID,
		Comment: input.Comment,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}

// DeleteTask deletes a task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Task not found"}})
		return
	}

	// Authorization check
	authUserID := uint(c.MustGet("user_id").(float64))
	if task.CreatedBy != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You are not authorized to delete this task"}})
		return
	}

	// Using a transaction to ensure atomicity
	tx := database.DB.Begin()

	// Delete associations
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.AssignTaskToUser{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.AssignTaskToGroup{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskFollowupUser{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskStatusUpdateLog{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskCommentLog{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskSeenByUser{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task associations"})
		return
	}

	// Delete the task itself
	if err := tx.Delete(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
