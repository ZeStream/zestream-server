package utils

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

/*
VideoIDGen returns an unique videoID and appends the fileExtension to it,
it takes the fileExtensionas parameter
*/
func VideoIDGen(fileExtension string) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a 8 digit random number
	randomNumber := rand.Intn(100000000) + 10000000

	// VideoID = (8-digit random number) + (file Name)
	return strconv.Itoa(randomNumber) + fileExtension
}

// WrapStringInQuotes returns the string wrapped in quotes
func WrapStringInQuotes(str string) string {
	var buff bytes.Buffer

	buff.WriteString(str)
	buff.WriteString(" ")

	return buff.String()
}
