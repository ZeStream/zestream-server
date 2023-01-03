package service

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"sync"
	"zestream-server/constants"
	"zestream-server/utils"

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

	var wg sync.WaitGroup

	wg.Add(len(constants.AudioFileTypeMap) + len(constants.VideoFileTypeMap))

	go generateAudioFiles(fileName, targetFile, outputPath, &wg)

	go generateVideoFiles(fileName, targetFile, outputPath, &wg)

	wg.Wait()

	generateMPD(outputPath)
}

func generateAudioFiles(_ string, targetFile string, outputPath string, wg *sync.WaitGroup) {
	for fileType, filePrefix := range constants.AudioFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateCappedBitrateAudio(targetFile, outputFile, fileType, wg)
	}
}

func generateVideoFiles(_ string, targetFile string, outputPath string, wg *sync.WaitGroup) {
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

func generateMPD(outputPath string) {
	var fileArgs bytes.Buffer

	mergePathToString(&fileArgs, outputPath, constants.AudioFileTypeMap)
	mergePathToString(&fileArgs, outputPath, constants.VideoFileTypeMap)

	var filePaths = strings.TrimSuffix(fileArgs.String(), " ")

	var inputArgsMap = map[string]string{
		constants.Mp4BoxArgs[constants.Dash]:        constants.Mp4BoxConfig[constants.Dash],
		constants.Mp4BoxArgs[constants.Rap]:         constants.Mp4BoxConfig[constants.Rap],
		constants.Mp4BoxArgs[constants.FragRap]:     constants.Mp4BoxConfig[constants.FragRap],
		constants.Mp4BoxArgs[constants.BsSwitching]: constants.Mp4BoxConfig[constants.BsSwitching],
		constants.Mp4BoxArgs[constants.Profile]:     constants.Mp4BoxConfig[constants.Profile],
		constants.Mp4BoxArgs[constants.Out]:         outputPath + constants.DashOutputExt,
	}

	inputArgsStr := utils.StringToArgsGenerator(inputArgsMap)

	var argsArr = strings.Split(inputArgsStr+filePaths, " ")

	cmd := exec.Command(constants.MP4Box, argsArr...)

	o, err := cmd.CombinedOutput()

	utils.DeleteFiles(filePaths)

	if err != nil {
		log.Println(err)
	}

	log.Println(string(o))
}

func mergePathToString(fileArgs *bytes.Buffer, outputPath string, fileTypes map[constants.FILE_TYPE]string) {
	for _, filePrefix := range fileTypes {
		var outputFile = outputPath + filePrefix
		if utils.IsFileValid(outputFile) {
			fileArgs.WriteString(utils.WrapStringInQuotes(outputFile))
		}
	}
}
