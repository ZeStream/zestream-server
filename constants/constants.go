package constants

const PORT = "PORT"
const DEFAULT_PORT = "8080"

const DOWNLOAD_FILE_PATH_PREFIX = "downloads"
const DOWNLOAD_FOLDER_PERM = 0666

const OUTPUT_FILE_PATH_PREFIX = "output"

type FILE_TYPE int

const (
	Audio192K FILE_TYPE = iota
	Video5M
	Video3M
	Video1M
	Video800K
	Video400K
)

var AudioFileType = map[FILE_TYPE]string{
	Audio192K: "_audio192k.m4a",
}

var VideoFileType = map[FILE_TYPE]string{
	Video5M:   "_5000.mp4",
	Video3M:   "_3000.mp4",
	Video1M:   "_1500.mp4",
	Video800K: "_800.mp4",
	Video400K: "_400.mp4",
}

var AudioBitrate = map[FILE_TYPE]string{
	Audio192K: "192k",
}

var VideoBitrate = map[FILE_TYPE]string{
	Video5M:   "5000k",
	Video3M:   "3000k",
	Video1M:   "1500k",
	Video800K: "800k",
	Video400K: "400k",
}

var VideoBufferSize = map[FILE_TYPE]string{
	Video5M:   "12000k",
	Video3M:   "6000k",
	Video1M:   "3000k",
	Video800K: "2000k",
	Video400K: "1000k",
}
