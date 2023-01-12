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
const DOWNLOAD_FILE_PATH_PREFIX = "downloads"
const OUTPUT_FILE_PATH_PREFIX = "output"
const CLOUD_CONTAINER_NAME = "zestream-dash"
const DOWNLOAD_FOLDER_PERM = 0666

const AWS_ENDPOINT = "http://localhost:4566"
const PRESIGNED_URL_EXPIRATION = 60 * time.Minute

// event queue
const RABBIT_MQ_CHANNEL = "VideoProcessing"
const RABBIT_MQ_TIMEOUT = 5 * time.Second
