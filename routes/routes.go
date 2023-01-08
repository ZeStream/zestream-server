package routes

import (
	"zestream-server/constants"
	"zestream-server/controllers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Init function will perform all route operations
func Init() *gin.Engine {

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	// middleware attaching a new transaction id for every request to context
	// transaction id will help to track all logs releated to any request better
	r.Use(func(ctx *gin.Context) {
		transactionID := ctx.Request.Header.Get("transaction-id")

		if transactionID == "" {
			transactionID = uuid.New().String()
		}

		ctx.Set(constants.TRANSACTION_ID_KEY, transactionID)
		ctx.Next()
	})

	// Create a new session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.S3_REGION),
	})
	if err != nil {
		return nil
	}

	// Create a new S3 client
	s3Client := s3.New(sess)

	v1 := r.Group("/api/v1")

	v1.GET("ping", controllers.Ping)

	v1.POST("process-video", controllers.ProcessVideo)

	v1.GET("generate-presigned-url", func(c *gin.Context) {
		controllers.GeneratePresignedURL(c, s3Client)
	})

	v1.POST("register_video_process", controllers.PublishMessage)

	return r
}
