package controllers

import (
	"fmt"
	"net/http"

	"taskmanager/database"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
)

type AssignUsersToGroupInput struct {
	UserIDs []uint `json:"user_ids" binding:"required,min=1"`
	GroupID uint   `json:"group_id" binding:"required"`
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
// AssignUsersToGroup assigns multiple users to a group
func AssignUsersToGroup(c *gin.Context) {
	var input AssignUsersToGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	// Check if GroupID exists
	var group models.Group
	if err := database.DB.First(&group, input.GroupID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Group not found"}})
		return
	}

	var failedAssignments []string
	for _, userID := range input.UserIDs {
		// Check if UserID exists
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			failedAssignments = append(failedAssignments, fmt.Sprintf("User with ID %d not found", userID))
			continue
		}

		// Check if the user is already in the group
		var existingUserGroup models.UserGroup
		if database.DB.Where("user_id = ? AND group_id = ?", userID, input.GroupID).First(&existingUserGroup).Error == nil {
			failedAssignments = append(failedAssignments, fmt.Sprintf("User %s is already in this group", user.Username))
			continue
		}
		
		userGroup := models.UserGroup{UserID: userID, GroupID: input.GroupID}
		if err := database.DB.Create(&userGroup).Error; err != nil {
			failedAssignments = append(failedAssignments, fmt.Sprintf("Failed to add user %s to group: %v", user.Username, err))
			continue
		}
	}

	if len(failedAssignments) > 0 {
		c.JSON(http.StatusMultiStatus, gin.H{"message": "Some users could not be assigned", "failed_assignments": failedAssignments})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Users assigned successfully"})
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
