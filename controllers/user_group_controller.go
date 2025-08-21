package controllers

import (
	"fmt"
	"net/http"

	"taskmanager/database"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
)

type AddUserToGroupInput struct {
	UserID  uint `json:"user_id" binding:"required"`
	GroupID uint `json:"group_id" binding:"required"`
}

// GetGroupsCreatedByUser list group created by auth user
func GetGroupsCreatedByUser(c *gin.Context) {
	userID := c.MustGet("user_id").(float64) // user_id is stored as float64 by jwt-go
	var groups []models.Group
	if err := database.DB.Preload("Users").Where("created_by = ?", uint(userID)).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups created by user"})
		return
	}

	var groupResponses []models.GroupResponse
	for _, group := range groups {
		var userResponses []models.UserResponse
		for _, user := range group.Users {
			var userGroup models.UserGroup
			if err := database.DB.Where("user_id = ? AND group_id = ?", user.ID, group.ID).First(&userGroup).Error; err != nil {
				// Handle error or skip if association not found (shouldn't happen if Preload is correct)
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

// AddUserToGroup adds a user to a group
func AddUserToGroup(c *gin.Context) {
	var input AddUserToGroupInput
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

	// Check if GroupID exists
	var group models.Group
	if err := database.DB.First(&group, input.GroupID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Group not found"}})
		return
	}

	// Check if the user is already in the group
	var existingUserGroup models.UserGroup
	if database.DB.Where("user_id = ? AND group_id = ?", input.UserID, input.GroupID).First(&existingUserGroup).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"errors": []string{"User is already in this group"}})
		return
	}
	
	userGroup := models.UserGroup{UserID: input.UserID, GroupID: input.GroupID}

	if err := database.DB.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to add user to group"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": userGroup})
}

// RemoveUserFromGroup removes a user from a group
func RemoveUserFromGroup(c *gin.Context) {
	id := c.Param("id")
	var userGroup models.UserGroup

	// 1. Check if the user-group association exists
	if err := database.DB.Where("id = ?", id).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User-group association not found"})
		return
	}

	// 2. Get the authenticated user's ID
	authUserID := uint(c.MustGet("user_id").(float64))

	// 3. Find the associated group
	var group models.Group
	if err := database.DB.First(&group, userGroup.GroupID).Error; err != nil {
		// This case should ideally not happen if data integrity is maintained
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Associated group not found"})
		return
	}

	// 4. Check if the authenticated user is the owner of the group
	if group.CreatedBy != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to remove users from this group"})
		return
	}

	// 5. Delete the user-group association
	database.DB.Delete(&userGroup)

	c.JSON(http.StatusOK, gin.H{"message": "User removed from group successfully"})
}

// GetGroupsForUser retrieves all groups the authenticated user belongs to
func GetGroupsForUser(c *gin.Context) {
	userID := c.MustGet("user_id").(float64) // user_id is stored as float64 by jwt-go
	var userGroups []models.UserGroup
	fmt.Println(userID)
	if err := database.DB.Where("user_id = ?", uint(userID)).Find(&userGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userGroups})
}
