package configs

import (
	"log"
	"os"

	dotEnv "github.com/joho/godotenv"
)

type CONFIG_KEY int

const (
	PORT CONFIG_KEY = iota
	RABBITMQ_URI
	AZURE_ACCOUNT_NAME
	AZURE_ACCESS_KEY
	AZURE_ENDPOINT
	AWS_S3_BUCKET_NAME
	AWS_S3_REGION
	GCP_BUCKET_NAME
	GCP_PROJECT_ID
)

var configVars = map[CONFIG_KEY]string{
	PORT:               "PORT",
	RABBITMQ_URI:       "RABBITMQ_URI",
	AZURE_ACCOUNT_NAME: "AZURE_ACCOUNT_NAME",
	AZURE_ACCESS_KEY:   "AZURE_ACCESS_KEY",
	AZURE_ENDPOINT:     "AZURE_ENDPOINT",
	AWS_S3_BUCKET_NAME: "AWS_S3_BUCKET_NAME",
	AWS_S3_REGION:      "AWS_S3_REGION",
	GCP_BUCKET_NAME:    "GCP_BUCKET_NAME",
	GCP_PROJECT_ID:     "GCP_PROJECT_ID",
}

var EnvVar = map[CONFIG_KEY]string{
	PORT:               "4444",
	RABBITMQ_URI:       "",
	AZURE_ACCOUNT_NAME: "",
	AZURE_ACCESS_KEY:   "",
	AZURE_ENDPOINT:     "",
	AWS_S3_BUCKET_NAME: "",
	AWS_S3_REGION:      "",
	GCP_BUCKET_NAME:    "",
	GCP_PROJECT_ID:     "",
}

func LoadEnv() {
	err := dotEnv.Load()

	if err != nil {
		log.Fatalln(err)
	}

	for key, val := range configVars {
		EnvVar[key] = getEnv(val, EnvVar[key])
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
