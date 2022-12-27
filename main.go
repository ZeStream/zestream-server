package main

import (
	"fmt"
	"net/http"
	"os"
	"zestream/constants"
	"zestream/routes"

	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	kafkaUriVal, isKafkaUriSet := os.LookupEnv("KAFKA_URI")
	if isKafkaUriSet == true && kafkaUriVal != "" {
		kafkaUri := os.Getenv(constants.KAFKA_URI)
	} else {
		fmt.Println("Kafka URI is not set")
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
