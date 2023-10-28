package controllers

import (
	"net/http"
	"path/filepath"
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

func GetPresignedURL(c *gin.Context) {
	fileName := c.Query("fileName")
	basePath := c.Query("basePath")
	extension := filepath.Ext(fileName)

	if fileName == "" || extension == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter 'fileName' or it has no extension"})
		return
	}

	videoID := utils.VideoIDGen(extension)

	url := utils.GetSignedURL(videoID, basePath)
	if url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL", "details": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preSignedURL": url, "videoID": videoID})
}
