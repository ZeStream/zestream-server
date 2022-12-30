package configs

import (
	"log"
	"os"

	dotEnv "github.com/joho/godotenv"
)

type CONFIG_KEY int

const (
	KAFKA_URI CONFIG_KEY = iota
	PORT
)

var configVars = map[CONFIG_KEY]string{
	PORT:      "PORT",
	KAFKA_URI: "KAFKA_URI",
}

var EnvVar = map[CONFIG_KEY]string{
	PORT:      "4444",
	KAFKA_URI: "",
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
