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
const READ_TIMEOUT = 5 * time.Second
const WRITE_TIMEOUT = 5 * time.Second
const IDLE_TIMEOUT = 30 * time.Second
const DEFAULT_THUMBNAIL_TIMESTAMP = "00:00:02"
const GCP_BUCKET_NAME = "zstream-bucket"
const GCP_PROJECT_ID = ""
const AZURE_ACCOUNT_NAME = ""
const AZURE_ENDPOINT = ""
const TRANSACTION_ID_KEY = "TRANSACTION_ID"
