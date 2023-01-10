package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/routes"
)

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

	err = r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
	}
}
