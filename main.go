package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/routes"
	"zestream-server/service"
	"zestream-server/utils"
)

func dev() {
	utils.Fetch("https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4", "Test.mp4")
	service.GenerateDash("Test.mp4")
}

func main() {
	configs.LoadEnv()

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	kafkaURI := os.Getenv("KAFKA_URI")
	if kafkaURI == "" {
		log.Fatal("Error: KAFKA_URI environment variable not set")
	}

	if port == "" {
		port = constants.DEFAULT_PORT
	}

	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Println(err)
	}

	r.Run(":" + port)
}
