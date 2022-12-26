package service

import (
	"log"
	"sync"
	"zestream/constants"
	"zestream/utils"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GenerateDash(fileName string) {
	targetFile, err := utils.GetDownloadFilePathName(fileName)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(constants.AudioFileType) + len(constants.VideoFileType))

	for fileType, filePrefix := range constants.AudioFileType {
		outputPath, err := utils.GetOutputFilePathName(fileName, "www")

		if err != nil {
			log.Println(err)
			return
		}

		var outputFile = outputPath + filePrefix

		go generateCappedBitrateAudio(targetFile, outputFile, fileType, &wg)
	}

	for fileType, filePrefix := range constants.VideoFileType {
		outputPath, err := utils.GetOutputFilePathName(fileName, "www")

		if err != nil {
			log.Println(err)
			return
		}

		var outputFile = outputPath + filePrefix

		go generateCappedBitrateVideo(targetFile, outputFile, fileType, &wg)
	}

	wg.Wait()
}

func generateCappedBitrateAudio(targetFile string, outputFile string, fileType constants.FILE_TYPE, wg *sync.WaitGroup) {
	ffmpeg.Input(targetFile).
		Output(outputFile, ffmpeg.KwArgs{
			"c:a":      "aac",
			"b:a":      constants.AudioBitrate[fileType],
			"allow_sw": 1,
			"vn":       "",
		}).
		OverWriteOutput().ErrorToStdOut().Run()

	wg.Done()
}

func generateCappedBitrateVideo(targetFile string, outputFile string, fileType constants.FILE_TYPE, wg *sync.WaitGroup) {
	ffmpeg.Input(targetFile).
		Output(outputFile, ffmpeg.KwArgs{
			"preset":   "fast",
			"tune":     "film",
			"fps_mode": "passthrough",
			"an":       "",
			"c:v":      "libx264",
			"crf":      "23",
			"maxrate":  constants.VideoBitrate[fileType],
			"bufsize":  constants.VideoBufferSize[fileType],
			"f":        "mp4",
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	wg.Done()
}
