package main

import (
	"fmt"
	"log"
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
