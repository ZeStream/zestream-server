package routes

import (
	"log"
	"zestream-server/controllers"

	"github.com/gin-gonic/gin"
)

// Init function will perform all route operations
func Init() *gin.Engine {
	log.Println("Running ZeStream as HTTP Server")
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	apiV1 := r.Group("/api/v1")

	r.GET("health", controllers.Ping)

	process := new(controllers.Process)

	// /api/v1
	apiV1.POST("video/process", process.Video)
	apiV1.GET("url/presigned", controllers.GetPresignedURL)

	// audio
	apiV1.POST("audio/process", process.AudioController)

	return r
}
