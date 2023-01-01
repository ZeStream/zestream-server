package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/docs"
	"zestream-server/routes"
	"zestream-server/service"
	"zestream-server/utils"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func dev() {
	utils.Fetch("https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4", "Test.mp4")
	service.GenerateDash("Test.mp4")
}

func main() {
	configs.LoadEnv()

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "ZeStream - An adaptive video streaming server"
	docs.SwaggerInfo.Description = "ZeStream is the backend service which you can self-deploy, and use its API to process the video and store it on a storage bucket like AWS S3/Google Cloud/Azure...."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ":" + port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}

	kafkaURI := os.Getenv("KAFKA_URI")
	if kafkaURI == "" {
		log.Fatal("Error: KAFKA_URI environment variable not set")
	}

	if port == "" {
		port = constants.DEFAULT_PORT
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println(err)
	}

	r.Run(":" + port)
}
