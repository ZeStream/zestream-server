package service

import (
	"log"
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

	for _, filePrefix := range constants.AudioFileType {
		outputPath, _ := utils.GetOutputFilePathName(fileName, "earth")

		ffmpeg.Input(targetFile).
			Output(outputPath+filePrefix, ffmpeg.KwArgs{"c:a": "aac", "b:a": "192k", "allow_sw": 1}).
			OverWriteOutput().ErrorToStdOut().Run()
	}

}
