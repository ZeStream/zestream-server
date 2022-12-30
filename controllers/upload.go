package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"zestream-server/constants"
	"zestream-server/types"
	"zestream-server/utils"
)

func GeneratePresignedURL(c *gin.Context) {

	// Validate the incoming request
	var request types.GenerateURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request parameters", "details": err.Error()})
		return
	}

	// Create a new session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(constants.S3_REGION),
	})
	svc := s3.New(sess)

	// Generate a New Video ID
	videoID := utils.VideoIDGen(request.FileName)

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

	c.JSON(http.StatusOK, gin.H{"pre-signed URL": urlStr, "videoID": videoID})
}
