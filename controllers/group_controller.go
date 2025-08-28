package controllers

import (
	"net/http"

	"taskmanager/database"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type CreateGroupInput struct {
	Label string `json:"label" binding:"required,min=3,max=100"`
}

type UpdateGroupInput struct {
	Label string `json:"label"`
}

// CreateGroup creates a new group
func CreateGroup(c *gin.Context) {
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			_ = ve.Translate(trans)

			for _, e := range ve {
				if e.Field() == "Label" && e.Tag() == "required" {
					errors = append(errors, "Group label is required")
				} else if e.Field() == "Label" && e.Tag() == "min" {
					errors = append(errors, "Group label must be at least 3 characters long")
				} else if e.Field() == "Label" && e.Tag() == "max" {
					errors = append(errors, "Group label cannot exceed 100 characters")
				} else {
					errors = append(errors, e.Translate(trans))
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// In a real application, you would get the CreatedBy from the authenticated user's ID
	// For now, we'll use a placeholder or assume it's handled by middleware
	group := models.Group{Label: input.Label, CreatedBy: uint(c.MustGet("user_id").(float64))}

	if err := database.DB.Create(&group).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"errors": []string{"Group with this label already exists"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to create group"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": group})
}

// GetGroups retrieves all groups
func GetGroups(c *gin.Context) {
	var groups []models.Group
	// Preload the Users relationship
	if err := database.DB.Preload("Users").Order("created_at DESC").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	var groupResponses []models.GroupResponse
	for _, group := range groups {
		var userResponses []models.UserResponse
		for _, user := range group.Users {
			var userGroup models.UserGroup
			if err := database.DB.Where("user_id = ? AND group_id = ?", user.ID, group.ID).First(&userGroup).Error; err != nil {
				// This should ideally not happen if the data is consistent, but handle it gracefully
				continue
			}
			userResponses = append(userResponses, models.UserResponse{
				ID:            user.ID,
				Username:      user.Username,
				Status:        user.Status,
				UserLabel:     user.UserLabel,
				AssociationID: userGroup.ID,
			})
		}
		groupResponses = append(groupResponses, models.GroupResponse{
			ID:        group.ID,
			Label:     group.Label,
			CreatedBy: group.CreatedBy,
			CreatedAt: group.CreatedAt,
			UpdatedAt: group.UpdatedAt,
			Users:     userResponses,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": groupResponses})
}

// GetGroupByID retrieves a single group by ID
func GetGroupByID(c *gin.Context) {
	id := c.Param("id")
	var group models.Group

	if err := database.DB.Preload("Users").Where("id = ?", id).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	var userResponses []models.UserResponse
	for _, user := range group.Users {
		userResponses = append(userResponses, models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Status:    user.Status,
			UserLabel: user.UserLabel,
		})
	}

	groupResponse := models.GroupResponse{
		ID:        group.ID,
		Label:     group.Label,
		CreatedBy: group.CreatedBy,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		Users:     userResponses,
	}

	c.JSON(http.StatusOK, gin.H{"data": groupResponse})
}

// UpdateGroup updates an existing group
func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group

	if err := database.DB.Where("id = ?", id).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var input UpdateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&group).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": group})
}

// DeleteGroup deletes a group
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group

	if err := database.DB.Where("id = ? ", id).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Authorization check: Only the creator can delete the group
	authUserID := uint(c.MustGet("user_id").(float64))
	if group.CreatedBy != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You are not authorized to delete this group"}})
		return
	}

	tx := database.DB.Begin()

	// Delete associated UserGroup entries
	if err := tx.Where("group_id = ?", group.ID).Delete(&models.UserGroup{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated user groups"})
		return
	}

	// Delete the group itself
	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
