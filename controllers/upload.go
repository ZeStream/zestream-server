package controllers

import (
	"net/http"
	"path/filepath"
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

func GeneratePresignedURL(c *gin.Context) {

	// Obtain the file name and extension from query params
	fileName := c.Query("fileName")
	extension := filepath.Ext(fileName)

	if fileName == "" || extension == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter 'fileName' or 'fileName' provided without any extension"})
		return
	}

	// Generate a New Video ID
	videoID := utils.VideoIDGen(extension)

	// Create a PutObjectRequest with the necessary parameters
	preSignedURL := utils.GetPreSignedURL(videoID)

	c.JSON(http.StatusOK, gin.H{"preSignedURL": preSignedURL, "videoID": videoID})
}
