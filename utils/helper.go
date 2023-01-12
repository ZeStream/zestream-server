package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	math_rand "math/rand"
	"strconv"
)

/*
VideoIDGen returns an unique videoID and appends the fileExtension to it,
it takes the fileExtensionas parameter
*/
func VideoIDGen(fileExtension string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return err.Error()
	}

	var i int64 = int64(binary.LittleEndian.Uint64(b[:]))
	math_rand.Seed(i)

	// Generate a 8 digit random number
	randomNumber := math_rand.Intn(99999999-10000000) + 10000000

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

// LogErr logs the given error
func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
