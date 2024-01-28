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

	process := new(controllers.Process)
	image := new(controllers.Image)

	apiV1 := r.Group("/api/v1")

	r.GET("health", controllers.Ping)

	// /api/v1
	apiV1.POST("video/process", process.Video)
	apiV1.GET("url/presigned", controllers.GetPresignedURL)

	r.GET("/i/*subpath", image.Get)

	return r
}
