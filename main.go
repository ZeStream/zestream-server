package main

import (
	"log"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/routes"
	"zestream-server/service"
)

func main() {
	configs.LoadEnv()

	r := routes.Init()

	port := os.Getenv(constants.PORT)

	// initialize RabbitMQ
	conn, ch, q, _, cancel := configs.InitRabbitMQ()

	configs.InitCloud()

	go service.VideoProcessConsumer(ch, q)

	err := r.Run(":" + port)
	failOnError(err)

	defer func() {
		defer conn.Close()
		defer ch.Close()
		defer cancel()
	}()
}

func failOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
