package routes

import (
	"zestream-server/controllers"

	"github.com/gin-gonic/gin"
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

	// TODO: write a functin to return session of AWS/GCP/Azure
	// Create a new session
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(constants.S3_REGION),
	// })
	// if err != nil {
	// 	return nil
	// }

	// // Create a new S3 client
	// s3Client := s3.New(sess)

	v1 := r.Group("/api/v1")

	v1.GET("ping", controllers.Ping)

	v1.POST("process-video", controllers.ProcessVideo)

	// v1.GET("generate-presigned-url", func(c *gin.Context) {
	// 	controllers.GeneratePresignedURL(c, s3Client)
	// })

	return r
}
