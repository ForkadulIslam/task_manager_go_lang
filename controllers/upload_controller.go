package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadAttachment handles file uploads to the 'uploads' directory.
func UploadAttachment(c *gin.Context) {
	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("get form file error: %s", err.Error())})
		return
	}

	// Generate a unique filename to prevent collisions
	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	
	// Define the path to save the file
	// Assuming 'uploads' directory is at the root of the application
	savePath := filepath.Join("uploads", filename)

	// Save the file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("upload file error: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "path": savePath})
}
