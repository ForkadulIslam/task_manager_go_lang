package controllers

import (
	"fmt"
	"net/http"
	"time"

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
	FollowUpGroups   []uint      `json:"follow_up_groups" binding:"omitempty,dive,gt=0"`
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
	FollowUpGroups   []uint      `json:"follow_up_groups" binding:"omitempty,dive,gt=0"`
}

type UpdateTaskStatusInput struct {
	Status string `json:"status" binding:"required,oneof=Pending 'In Progress' 'In Review' Completed"`
}

type AddTaskCommentInput struct {
	Comment string `json:"comment" binding:"required"`
}

type GetMyTasksFilterInput struct {
	FromDate   *utils.Date `json:"from_date"`
	ToDate     *utils.Date `json:"to_date"`
	Status     string      `json:"status"`
	TaskTypeID uint        `json:"task_type_id"`
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

	// Validate FollowUpGroups
	for _, groupID := range input.FollowUpGroups {
		var group models.Group
		if err := database.DB.First(&group, groupID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Follow-up group with ID %d not found", groupID)}})
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

		// Create notification for assigned user
		var notification models.Notification
		notification = models.Notification{
			UserID:  userID,
			TaskID:  task.ID,
			Type:    "new_task",
			Message: fmt.Sprintf("You have been assigned a new task: %s", task.Label),
		}
		if err := database.DB.Create(&notification).Error; err != nil {
			// Handle error
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

		// Create notification for follow-up user
		notification := models.Notification{
			UserID:  userID,
			TaskID:  task.ID,
			Type:    "new_task",
			Message: fmt.Sprintf("You are following a new task: %s", task.Label),
		}
		if err := database.DB.Create(&notification).Error; err != nil {
			// Handle error
		}
	}

	// Assign task to follow-up groups
	for _, groupID := range input.FollowUpGroups {
		followUpGroup := models.TaskFollowupGroup{
			TaskID:  task.ID,
			GroupID: groupID,
		}
		if err := database.DB.Create(&followUpGroup).Error; err != nil {
			// Handle error, perhaps rollback task creation or log it
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to follow-up group"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

// GetTasks retrieves all tasks created by the authenticated user
func GetTasks(c *gin.Context) {
	authUserID := uint(c.MustGet("user_id").(float64))

	var tasks []models.Task
	if err := database.DB.
		Preload("AssignedUsers.User").
		Preload("AssignedGroups.Group.Users").
		Preload("FollowupUsers.User").
		Preload("FollowupGroups.Group.Users").
		Where("created_by = ?", authUserID).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTaskByID retrieves a single task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := database.DB.Preload("AssignedUsers.User").Preload("AssignedGroups.Group.Users").Preload("FollowupUsers.User").Preload("FollowupGroups.Group.Users").Preload("Comments.User").Preload("Creator").Where("id = ?", id).First(&task).Error; err != nil {
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

	// Validate FollowUpGroups, if provided
	if input.FollowUpGroups != nil {
		for _, groupID := range input.FollowUpGroups {
			var group models.Group
			if err := database.DB.First(&group, groupID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errors": []string{fmt.Sprintf("Follow-up group with ID %d not found", groupID)}})
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

	if input.FollowUpGroups != nil {
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskFollowupGroup{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update follow-up groups"})
			return
		}
		for _, groupID := range input.FollowUpGroups {
			followUpGroup := models.TaskFollowupGroup{TaskID: task.ID, GroupID: groupID}
			if err := tx.Create(&followUpGroup).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task to follow-up group"})
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

	// Create notification for task creator
	if task.CreatedBy != authUserID {
		var user models.User
		database.DB.First(&user, authUserID)
		notification := models.Notification{
			UserID:  task.CreatedBy,
			TaskID:  task.ID,
			Type:    "status_update",
			Message: fmt.Sprintf("Task '%s' status updated to '%s' by %s", task.Label, input.Status, user.Username),
		}
		if err := tx.Create(&notification).Error; err != nil {
			// Handle error
		}
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

	if followupAssignment > 0 {
		return true, nil
	}

	var followupGroupAssignment int64
	err = db.Model(&models.UserGroup{}).
		Joins("JOIN task_followup_groups ON user_groups.group_id = task_followup_groups.group_id").
		Where("task_followup_groups.task_id = ? AND user_groups.user_id = ?", taskID, userID).
		Count(&followupGroupAssignment).Error
	if err != nil {
		return false, err
	}

	return followupGroupAssignment > 0, nil

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

	// Allow task creator to comment
	if task.CreatedBy == authUserID {
		allowed = true
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

	// Create notifications for task creator and associated users
	var usersToNotify []uint
	usersToNotify = append(usersToNotify, task.CreatedBy)

	var assignedUsers []models.AssignTaskToUser
	database.DB.Where("task_id = ?", task.ID).Find(&assignedUsers)
	for _, u := range assignedUsers {
		usersToNotify = append(usersToNotify, u.UserID)
	}

	var followupUsers []models.TaskFollowupUser
	database.DB.Where("task_id = ?", task.ID).Find(&followupUsers)
	for _, u := range followupUsers {
		usersToNotify = append(usersToNotify, u.UserID)
	}

	// Remove duplicates and the user who commented
	userMap := make(map[uint]bool)
	for _, userID := range usersToNotify {
		if userID != authUserID {
			userMap[userID] = true
		}
	}

	var user models.User
	database.DB.First(&user, authUserID)

	for userID := range userMap {
		notification := models.Notification{
			UserID:  userID,
			TaskID:  task.ID,
			Type:    "new_comment",
			Message: fmt.Sprintf("New comment on task '%s' by %s", task.Label, user.Username),
		}
		if err := database.DB.Create(&notification).Error; err != nil {
			// Handle error
		}
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
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskFollowupGroup{}).Error; err != nil {
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

// GetMyTasks retrieves tasks assigned to or followed by the authenticated user
func GetMyTasks(c *gin.Context) {
	authUserID := uint(c.MustGet("user_id").(float64))

	// Get all group IDs for the current user
	var groupIDs []uint
	if err := database.DB.Model(&models.UserGroup{}).Where("user_id = ?", authUserID).Pluck("group_id", &groupIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user groups"})
		return
	}

	taskIDMap := make(map[uint]bool)

	// 1. Tasks assigned directly to the user
	var assignedUserTaskIDs []uint
	database.DB.Model(&models.AssignTaskToUser{}).Where("user_id = ?", authUserID).Pluck("task_id", &assignedUserTaskIDs)
	for _, id := range assignedUserTaskIDs {
		taskIDMap[id] = true
	}

	// 2. Tasks followed directly by the user
	var followupUserTaskIDs []uint
	database.DB.Model(&models.TaskFollowupUser{}).Where("user_id = ?", authUserID).Pluck("task_id", &followupUserTaskIDs)
	for _, id := range followupUserTaskIDs {
		taskIDMap[id] = true
	}

	if len(groupIDs) > 0 {
		// 3. Tasks assigned to the user's groups
		var assignedGroupTaskIDs []uint
		database.DB.Model(&models.AssignTaskToGroup{}).Where("group_id IN ?", groupIDs).Pluck("task_id", &assignedGroupTaskIDs)
		for _, id := range assignedGroupTaskIDs {
			taskIDMap[id] = true
		}

		// 4. Tasks followed by the user's groups
		var followupGroupTaskIDs []uint
		database.DB.Model(&models.TaskFollowupGroup{}).Where("group_id IN ?", groupIDs).Pluck("task_id", &followupGroupTaskIDs)
		for _, id := range followupGroupTaskIDs {
			taskIDMap[id] = true
		}
	}

	var relevantTaskIDs []uint
	for id := range taskIDMap {
		relevantTaskIDs = append(relevantTaskIDs, id)
	}

	var tasks []models.Task
	if len(relevantTaskIDs) > 0 {
		if err := database.DB.
			Preload("AssignedUsers.User").
			Preload("AssignedGroups.Group.Users").
			Preload("FollowupUsers.User").
			Preload("FollowupGroups.Group.Users").
			Preload("Creator").
			Where("id IN ?", relevantTaskIDs).
			Order("created_at DESC").
			Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
	} else {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetMyTasksFiltered retrieves tasks assigned to or followed by the authenticated user with filters
func GetMyTasksFiltered(c *gin.Context) {
	authUserID := uint(c.MustGet("user_id").(float64))

	var filterInput GetMyTasksFilterInput
	if err := c.ShouldBindJSON(&filterInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter input"})
		return
	}

	// Get all group IDs for the current user
	var groupIDs []uint
	if err := database.DB.Model(&models.UserGroup{}).Where("user_id = ?", authUserID).Pluck("group_id", &groupIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user groups"})
		return
	}

	taskIDMap := make(map[uint]bool)

	// 1. Tasks assigned directly to the user
	var assignedUserTaskIDs []uint
	database.DB.Model(&models.AssignTaskToUser{}).Where("user_id = ?", authUserID).Pluck("task_id", &assignedUserTaskIDs)
	for _, id := range assignedUserTaskIDs {
		taskIDMap[id] = true
	}

	// 2. Tasks followed directly by the user
	var followupUserTaskIDs []uint
	database.DB.Model(&models.TaskFollowupUser{}).Where("user_id = ?", authUserID).Pluck("task_id", &followupUserTaskIDs)
	for _, id := range followupUserTaskIDs {
		taskIDMap[id] = true
	}

	if len(groupIDs) > 0 {
		// 3. Tasks assigned to the user's groups
		var assignedGroupTaskIDs []uint
		database.DB.Model(&models.AssignTaskToGroup{}).Where("group_id IN ?", groupIDs).Pluck("task_id", &assignedGroupTaskIDs)
		for _, id := range assignedGroupTaskIDs {
			taskIDMap[id] = true
		}

		// 4. Tasks followed by the user's groups
		var followupGroupTaskIDs []uint
		database.DB.Model(&models.TaskFollowupGroup{}).Where("group_id IN ?", groupIDs).Pluck("task_id", &followupGroupTaskIDs)
		for _, id := range followupGroupTaskIDs {
			taskIDMap[id] = true
		}
	}

	var relevantTaskIDs []uint
	for id := range taskIDMap {
		relevantTaskIDs = append(relevantTaskIDs, id)
	}

	var tasks []models.Task
	if len(relevantTaskIDs) > 0 {
		db := database.DB.
			Preload("AssignedUsers.User").
			Preload("AssignedGroups.Group.Users").
			Preload("FollowupUsers.User").
			Preload("FollowupGroups.Group.Users").
			Preload("Creator").
			Where("id IN ?", relevantTaskIDs).
			Order("created_at DESC")

		// Apply filters from JSON body
		if filterInput.FromDate != nil && !filterInput.FromDate.IsZero() {
			db = db.Where("start_date >= ?", filterInput.FromDate.Time)
		}
		if filterInput.ToDate != nil && !filterInput.ToDate.IsZero() {
			// To include the entire end day, add 23 hours, 59 minutes, 59 seconds
			endOfDay := filterInput.ToDate.Time.Add(24 * time.Hour).Add(-time.Second)
			db = db.Where("due_date <= ?", endOfDay)
		}
		if filterInput.Status != "" {
			db = db.Where("status = ?", filterInput.Status)
		}
		if filterInput.TaskTypeID != 0 {
			db = db.Where("task_type_id = ?", filterInput.TaskTypeID)
		}

		if err := db.Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
	} else {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}
