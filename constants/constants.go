package constants

import (
	"time"
)

// server
const PORT = "PORT"
const DEFAULT_PORT = "8080"
const READ_TIMEOUT = 5 * time.Second
const WRITE_TIMEOUT = 5 * time.Second
const IDLE_TIMEOUT = 30 * time.Second

// folders
type FOLDER_TYPE int

const (
	Dashes FOLDER_TYPE = iota
	Images
	Audios
)

var CloudContainerNames = map[FOLDER_TYPE]string{
	Dashes: "dashes",
	Images: "i",
	Audios: "a",
}

const DOWNLOAD_FILE_PATH_PREFIX = "assets/downloads"
const OUTPUT_FILE_PATH_PREFIX = "assets/output"
const DOWNLOAD_FOLDER_PERM = 0666

const AWS_ENDPOINT = "http://localhost:4566"
const PRESIGNED_URL_EXPIRATION = 24 * time.Hour

// event queue
const RABBIT_MQ_CHANNEL = "VideoProcessing"
const RABBIT_MQ_TIMEOUT = 5 * time.Second
