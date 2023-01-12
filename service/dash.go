package service

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"sync"
	"zestream-server/constants"
	"zestream-server/types"
	"zestream-server/utils"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func generateDash(fileName string, watermark types.WaterMark) {
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

	wg.Add(len(constants.AudioFileTypeMap) + len(constants.VideoFileTypeMap) + 1)

	generateAudioFiles(targetFile, outputPath, &wg)

	generateVideoFiles(targetFile, outputPath, watermark, &wg)

	generateThumbnailFiles(targetFile, outputPath, &wg)

	wg.Wait()

	generateMPD(outputPath)

	deleteVideoFiles(outputPath)
	utils.DeleteFile(targetFile)
}

func generateAudioFiles(targetFile string, outputPath string, wg *sync.WaitGroup) {
	for fileType, filePrefix := range constants.AudioFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateMultiBitrateAudio(targetFile, outputFile, fileType, wg)
	}
}

func generateVideoFiles(targetFile string, outputPath string, watermark types.WaterMark, wg *sync.WaitGroup) {
	for fileType, filePrefix := range constants.VideoFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateMultiBitrateVideo(targetFile, outputFile, fileType, watermark, wg)
	}
}

func generateThumbnailFiles(targetFile string, outputPath string, wg *sync.WaitGroup) {
	for _, filePrefix := range constants.ImageFileTypeMap {
		var outputFile = outputPath + filePrefix

		go generateThumbnails(targetFile, outputFile, constants.DEFAULT_THUMBNAIL_TIMESTAMP, wg)
	}
}

func generateMultiBitrateAudio(targetFile string, outputFile string, fileType constants.FILE_TYPE, wg *sync.WaitGroup) {
	err := ffmpeg.Input(targetFile, ffmpeg.KwArgs{
		constants.AudioKwargs[constants.HWAccel]: constants.FFmpegConfig[constants.HWAccel],
	}).
		Output(outputFile, ffmpeg.KwArgs{
			constants.AudioKwargs[constants.AudioCodec]:        constants.FFmpegConfig[constants.AudioCodec],
			constants.AudioKwargs[constants.AudioBitrate]:      constants.AudioBitrateMap[fileType],
			constants.AudioKwargs[constants.AllowSoftEncoding]: constants.FFmpegConfig[constants.AllowSoftEncoding],
			constants.AudioKwargs[constants.VideoExclusion]:    constants.FFmpegConfig[constants.VideoExclusion],
		}).
		OverWriteOutput().Run()
	if err != nil {
		m := "error generating multi bitrate audio"
		log.Println(m, err)
	}

	wg.Done()
}

func generateMultiBitrateVideo(targetFile string, outputFile string, fileType constants.FILE_TYPE, watermark types.WaterMark, wg *sync.WaitGroup) {

	err := getInput(targetFile, watermark).
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
		Run()

	if err != nil {
		m := "error generating multi bitrate video"
		log.Println(m, err)
	}

	wg.Done()
}

// generateThumbnail generates a thumbnail at given timestamp, from the target file and write it to output file
func generateThumbnails(targetFile string, outputFile string, timeStamp string, wg *sync.WaitGroup) {
	err := ffmpeg.Input(targetFile).
		Output(outputFile, ffmpeg.KwArgs{
			constants.VideoKwargs[constants.ScreenShot]:  timeStamp,
			constants.VideoKwargs[constants.VideoFrames]: constants.FFmpegConfig[constants.VideoFrames],
		}).
		OverWriteOutput().
		Run()

	if err != nil {
		m := "error generating thumbnail"
		log.Println(m, err)
	}

	wg.Done()
}

// generateMPD (Media Presentation Description), generates the XML description file
// for the given output path.
func generateMPD(outputPath string) {
	var fileArgs bytes.Buffer

	utils.CheckFileExistsAndAppendToBuffer(&fileArgs, outputPath, constants.AudioFileTypeMap)
	utils.CheckFileExistsAndAppendToBuffer(&fileArgs, outputPath, constants.VideoFileTypeMap)

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

	if err != nil {
		log.Println(err)
	}

	log.Println(string(o))
}

// deleteVideoFiles deletes the videos given in the path
func deleteVideoFiles(outputPath string) {
	for _, filePrefix := range constants.VideoFileTypeMap {
		var file = outputPath + filePrefix

		utils.DeleteFile(file)
	}
}

// Returns input based on if the watermark is needed or not
func getInput(targetFile string, watermark types.WaterMark) *ffmpeg.Stream {

	input := ffmpeg.Input(targetFile)

	if !watermark.IsEmpty() {
		watermarkPath, err := utils.GetDownloadFilePathName(watermark.ID)
		if err != nil {
			log.Println(err)
		}

		filterArgs := "" + watermark.Position.X + ":" + watermark.Position.Y + ""

		return ffmpeg.Filter(
			[]*ffmpeg.Stream{
				input,
				getOverlay(watermarkPath, watermark),
			}, constants.Overlay, ffmpeg.Args{filterArgs})
	}
	return input
}

// Returns an overlay of watermark which can be used in Filter
func getOverlay(watermarkPath string, watermark types.WaterMark) *ffmpeg.Stream {
	overlayArgs := "" + watermark.Dimension.X + ":" + watermark.Dimension.Y + ""
	return ffmpeg.Input(watermarkPath).Filter(constants.Scale, ffmpeg.Args{overlayArgs})
}
