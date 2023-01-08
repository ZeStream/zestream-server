package main

import (
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

	conn, ch, q, _, cancel := configs.InitRabbitMQ()

	go service.VideoProcessConsumer(ch, q)

	r.Run(":" + port)

	defer func() {
		defer conn.Close()
		defer ch.Close()
		defer cancel()
	}()
}
