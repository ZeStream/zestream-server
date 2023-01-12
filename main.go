package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/logger"
	"zestream-server/routes"
)

func main() {
	configs.LoadEnv()

	logger.Init("zestream-server", os.Getenv("LOG_LEVEL"))

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	kafkaURI := os.Getenv("KAFKA_URI")
	if kafkaURI == "" {
		logger.Error(context.TODO(), "KAFKA_URI environment variable not set", logger.Z{
			"kafka_uri": os.Getenv("KAFKA_URI"),
		})
		panic("KAFKA_URI environment variable not set")
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

	err = r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
	}
}
