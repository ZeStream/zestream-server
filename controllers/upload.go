package controllers

import (
	"net/http"
	"path/filepath"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"
)

func GeneratePresignedAWSURL(c *gin.Context, s3Client s3iface.S3API) {

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
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(configs.EnvVar[configs.AWS_S3_BUCKET_NAME]),
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
