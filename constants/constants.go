package constants

import (
	"time"
)

const PORT = "PORT"
const DEFAULT_PORT = "8080"

const DOWNLOAD_FILE_PATH_PREFIX = "downloads"
const DOWNLOAD_FOLDER_PERM = 0666
const S3_BUCKET_NAME = "zstream-bucket"
const S3_REGION = "us-east-1"
const AWS_ENDPOINT = "http://localhost:4566"
const PRESIGNED_URL_EXPIRATION = 60 * time.Minute
const OUTPUT_FILE_PATH_PREFIX = "output"
const DEFAULT_THUMBNAIL_TIMESTAMP = "00:00:02"
