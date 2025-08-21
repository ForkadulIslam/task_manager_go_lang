package controllers

import (
	"net/http"

	"taskmanager/database"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
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
	group := models.Group{Label: input.Label, CreatedBy: 1} // Placeholder CreatedBy

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
	if err := database.DB.Preload("Users").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	var groupResponses []models.GroupResponse
	for _, group := range groups {
		var userResponses []models.UserResponse
		for _, user := range group.Users {
			userResponses = append(userResponses, models.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Status:    user.Status,
				UserLabel: user.UserLabel,
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

	if err := database.DB.Where("id = ?", id).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": group})
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

	if err := database.DB.Where("id = ?", id).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	database.DB.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}


