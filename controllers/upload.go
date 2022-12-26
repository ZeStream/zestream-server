package controllers

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func UploadFile(c *gin.Context) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error uploading file"})
		return
	}

	// Upload the file to S3.
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error opening file"})
		return
	}

	defer fileBytes.Close()
	fileData, err := ioutil.ReadAll(fileBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file"})
		return
	}

	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(fileData),
		Bucket: aws.String("my-bucket"),
		Key:    aws.String(file.Filename),
	}

	_, err = svc.PutObject(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to S3"})
		return
	}

	// Get presigned URL for the uploaded file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("my-bucket"),
		Key:    aws.String(file.Filename),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": urlStr})

}

func GeneratePresignedURL(c *gin.Context) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	//The GetObjectRequest function is used to create a request for the S3 object for which presigned URL is to be generated
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("my-bucket"),
		Key:    aws.String("my-object-key"),
	})

	urlStr, err := req.Presign(60 * time.Minute)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": urlStr})

}
