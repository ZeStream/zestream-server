package utils

import (
	"os"
)

var Config = map[string]string{
	"KAFKA_URL": os.Getenv("KAFKA_URI"),
}
