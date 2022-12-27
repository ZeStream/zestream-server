package controllers

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
	"zestream-server/constants"
	"zestream-server/types"
)

func UploadFile(c *gin.Context) {

	// Creates a new session
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(constants.S3_REGION),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(constants.AWS_ENDPOINT),
	})
	svc := s3.New(sess)

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error uploading file", "details": err.Error()})
		return
	}

	// Upload the file to S3.
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error opening file", "details": err.Error()})
		return
	}

	// Convert Bytes into a File Object acceptable by s3.PutObjectInput
	defer func(fileBytes multipart.File) {
		err := fileBytes.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error closing file", "details": err.Error()})
			return
		}
	}(fileBytes)

	fileData, err := ioutil.ReadAll(fileBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file", "details": err.Error()})
		return
	}

	// Upload the file to S3.
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(fileData),
		Bucket: aws.String(constants.S3_BUCKET_NAME),
		Key:    aws.String(file.Filename),
	}

	_, err = svc.PutObject(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to S3", "details": err.Error()})
		return
	}

	// Get presigned URL for the uploaded file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.S3_BUCKET_NAME),
		Key:    aws.String(file.Filename),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": urlStr})

}

func GeneratePresignedURL(c *gin.Context) {

	// Validate the incoming request
	var request types.GenerateURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request parameters", "details": err.Error()})
		return
	}

	// Create a new session
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(constants.S3_REGION),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(constants.AWS_ENDPOINT),
	})
	svc := s3.New(sess)

	//The GetObjectRequest function is used to create a request for the S3 object for which presigned URL is to be generated
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(constants.S3_BUCKET_NAME),
		Key:    aws.String(request.FileName),
	})

	//The Presign function is used to generate the presigned URL which is valid for the given time limit
	urlStr, err := req.Presign(constants.PRESIGNED_URL_EXPIRATION)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": urlStr})
}
