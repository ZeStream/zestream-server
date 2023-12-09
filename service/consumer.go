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

func ProcessConsumer(ch *rmq.Channel, q *rmq.Queue) {
	log.Println("Running ZeStream as Consumer")
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
		q.Name,            // queue
		"ProcessConsumer", // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	if err != nil {
		log.Panicln(err)
	}

	go func() {
		guard := make(chan int, maxProcesses)

		for d := range msgs {
			guard <- 1

			var video types.Video
			videoErr := json.Unmarshal(d.Body, &video)
			if videoErr != nil {
				log.Println(videoErr)
				continue
			}

			if video.Type == "mp4" {
				go processVideo(&video, guard)
				continue
			}

			var audio types.Audio
			audioErr := json.Unmarshal(d.Body, &audio)
			if err != nil {
				log.Println(audioErr)
				continue
			}

			if audio.Type == "mp3" {
				go processAudio(&audio, guard)
			}
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

	generateVideoDash(videoFileName, video.Watermark)

	uploader := utils.GetUploader(constants.CloudContainerNames[constants.Dashes], video.ID)

	outputDir, err := utils.GetOutputFilePathName(videoFileName, "")
	utils.LogErr(err)

	utils.UploadToCloudStorage(uploader, outputDir)

	err = os.RemoveAll(outputDir)
	utils.LogErr(err)

	<-guard
}

func processAudio(audio *types.Audio, guard <-chan int) {
	log.Println("Processing Audio: ", audio)

	var fileName = audio.ID + "." + audio.Type

	err := utils.Fetch(audio.Src, fileName)
	if err != nil {
		utils.LogErr(err)
		return
	}

	generateAudioDash(fileName)

	uploader := utils.GetUploader(constants.CloudContainerNames[constants.Dashes], audio.ID)

	outputDir, err := utils.GetOutputFilePathName(fileName, "")
	utils.LogErr(err)

	utils.UploadToCloudStorage(uploader, outputDir)

	err = os.RemoveAll(outputDir)
	utils.LogErr(err)

	<-guard
}
