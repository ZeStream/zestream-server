package service

import (
	"log"
	"sync"
	"time"
	"zestream/constants"
	"zestream/utils"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GenerateDash(fileName string) {
	targetFile, err := utils.GetDownloadFilePathName(fileName)
	if err != nil {
		log.Println(err)
	}

	var fileNameStripped = utils.RemoveExtensionFromFile(fileName)

	outputPath, err := utils.GetOutputFilePathName(fileName, fileNameStripped)
	if err != nil {
		log.Println(err)
		return
	}

	start := time.Now()

	var wg sync.WaitGroup

	wg.Add(len(constants.AudioFileTypeMap) + len(constants.VideoFileTypeMap))

	generateAudioFiles(fileName, targetFile, outputPath, &wg)

	generateVideoFiles(fileName, targetFile, outputPath, &wg)

	wg.Wait()

	log.Printf("main, execution time %s\n", time.Since(start))

}

func generateAudioFiles(fileName string, targetFile string, outputPath string, wg *sync.WaitGroup) {
	for fileType, filePrefix := range constants.AudioFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateCappedBitrateAudio(targetFile, outputFile, fileType, wg)
	}

}

func generateVideoFiles(fileName string, targetFile string, outputPath string, wg *sync.WaitGroup) {
	for fileType, filePrefix := range constants.VideoFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateCappedBitrateVideo(targetFile, outputFile, fileType, wg)
	}
}

func generateCappedBitrateAudio(targetFile string, outputFile string, fileType constants.FILE_TYPE, wg *sync.WaitGroup) {
	ffmpeg.Input(targetFile, ffmpeg.KwArgs{
		constants.AudioKwargs[constants.HWAccel]: constants.FFmpegConfig[constants.HWAccel],
	}).
		Output(outputFile, ffmpeg.KwArgs{
			constants.AudioKwargs[constants.AudioCodec]:        constants.FFmpegConfig[constants.AudioCodec],
			constants.AudioKwargs[constants.AudioBitrate]:      constants.AudioBitrateMap[fileType],
			constants.AudioKwargs[constants.AllowSoftEncoding]: constants.FFmpegConfig[constants.AllowSoftEncoding],
			constants.AudioKwargs[constants.VideoExclusion]:    constants.FFmpegConfig[constants.VideoExclusion],
		}).
		OverWriteOutput().ErrorToStdOut().Run()

	wg.Done()
}

func generateCappedBitrateVideo(targetFile string, outputFile string, fileType constants.FILE_TYPE, wg *sync.WaitGroup) {
	ffmpeg.Input(targetFile).
		Output(outputFile, ffmpeg.KwArgs{
			constants.VideoKwargs[constants.Preset]:         constants.FFmpegConfig[constants.Preset],
			constants.VideoKwargs[constants.Tune]:           constants.FFmpegConfig[constants.Tune],
			constants.VideoKwargs[constants.FpsMode]:        constants.FFmpegConfig[constants.FpsMode],
			constants.VideoKwargs[constants.AudioExclusion]: constants.FFmpegConfig[constants.AudioExclusion],
			constants.VideoKwargs[constants.VideoCodec]:     constants.FFmpegConfig[constants.VideoCodec],
			constants.VideoKwargs[constants.MaxRate]:        constants.VideoBitrateMap[fileType],
			constants.VideoKwargs[constants.BufferSize]:     constants.VideoBufferSizeMap[fileType],
			constants.VideoKwargs[constants.VideoFormat]:    constants.FFmpegConfig[constants.VideoFormat],
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	wg.Done()
}
