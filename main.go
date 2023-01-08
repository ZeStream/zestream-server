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
)

func dev() {
	// utils.Fetch("https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4", "Test.mp4")
	service.GenerateDash("Test.mp4", constants.WaterMark{FileName: "TestWatermark.png", Dimension: map[string]int{"x": 64, "y": -1}, Position: map[string]int{"x": 10, "y": 10}}, false)
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

	server := &http.Server{
		Addr:         port,
		ReadTimeout:  constants.READ_TIMEOUT,
		WriteTimeout: constants.WRITE_TIMEOUT,
		IdleTimeout:  constants.IDLE_TIMEOUT,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

	r.Run(":" + port)
}
