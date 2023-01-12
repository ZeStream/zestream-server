package service

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"zestream-server/configs"
	"zestream-server/constants"
	"zestream-server/types"
	"zestream-server/utils"

	rmq "github.com/rabbitmq/amqp091-go"
)

func VideoProcessConsumer(ch *rmq.Channel, q *rmq.Queue) {
	var forever chan struct{}

	maxProcesses, err := strconv.Atoi(configs.EnvVar[configs.MAX_CONCURRENT_PROCESSING])
	if err != nil {
		maxProcesses = 1
	}

	err = ch.Qos(
		maxProcesses, // prefetch count
		0,            // prefetch size
		false,        // global
	)
	utils.LogErr(err)

	msgs, err := ch.Consume(
		q.Name,                 // queue
		"VideoProcessConsumer", // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)

	if err != nil {
		log.Panicln(err)
	}

	go func() {
		guard := make(chan int, maxProcesses)

		for d := range msgs {
			guard <- 1

			var video types.Video

			log.Println("Request Consumed: ", video)

			err := json.Unmarshal(d.Body, &video)
			if err != nil {
				log.Println(err)
				continue
			}

			go processVideo(&video, guard)
		}
	}()

	<-forever
}

func processVideo(video *types.Video, guard <-chan int) {
	log.Println("Processing Video: ", video)

	var videoFileName = video.ID + "." + video.Type
	var waterMarkFileName = video.Watermark.ID + "." + video.Watermark.Type

	err := utils.Fetch(video.Src, videoFileName)
	if err != nil {
		utils.LogErr(err)
		return
	}

	if !video.Watermark.IsEmpty() {
		err = utils.Fetch(video.Watermark.Src, waterMarkFileName)
		utils.LogErr(err)
	}

	generateDash(videoFileName, video.Watermark)

	uploader := utils.GetUploader(constants.CLOUD_CONTAINER_NAME, video.ID)

	outputDir, err := utils.GetOutputFilePathName(videoFileName, "")
	utils.LogErr(err)

	err = os.RemoveAll(outputDir)
	utils.LogErr(err)

	utils.UploadToCloudStorage(uploader, outputDir)

	<-guard
}
