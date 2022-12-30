package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func VideoIDGen(fileExtension string) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a 8 digit random number
	randomNumber := rand.Intn(100000000) + 10000000

	// VideoID = (8-digit random number) + (file Name)
	return strconv.Itoa(randomNumber) + "-" + fileExtension
}
