package service

import (
	"encoding/json"
	"log"
	"zestream-server/types"
	"zestream-server/utils"

	rmq "github.com/rabbitmq/amqp091-go"
)

func VideoProcessConsumer(ch *rmq.Channel, q *rmq.Queue) {
	var forever chan struct{}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Panicln(err)
	}

	go func() {
		for d := range msgs {
			var video types.Video

			log.Println("Request Consumed: ", video)

			err := json.Unmarshal(d.Body, &video)
			if err != nil {
				log.Println(err)
				continue
			}

			processVideo(&video)
		}
	}()

	<-forever
}

func processVideo(video *types.Video) {
	log.Println("Processing Video: ", video)

	var fileName = video.ID + "." + video.Type

	utils.Fetch(video.Src, fileName)

	generateDash(fileName, video.Watermark)
}
