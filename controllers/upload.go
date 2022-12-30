package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"zestream-server/constants"
	"zestream-server/utils"
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

	// Create a new session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(constants.S3_REGION),
	})
	svc := s3.New(sess)

	// Create a PutObjectRequest with the necessary parameters
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(constants.S3_BUCKET_NAME),
		Key:    aws.String(videoID),
	})

	//  Sign the request and generate a presigned URL
	urlStr, err := req.Presign(constants.PRESIGNED_URL_EXPIRATION)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preSignedURL": urlStr, "videoID": videoID})
}
