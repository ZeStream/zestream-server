package main

import (
	"flag"
	"log"
	"os"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/routes"
	"zestream-server/service"
)

func main() {
	isConsumer := flag.Bool("consumer", false, "run service as consumer")
	flag.Parse()

	configs.LoadEnv()
	configs.InitCloud()
	conn, ch, q, _, cancel := configs.InitRabbitMQ()

	if *isConsumer {
		service.ProcessConsumer(ch, q)
		return
	}

	r := routes.Init()
	port := os.Getenv(constants.PORT)

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
