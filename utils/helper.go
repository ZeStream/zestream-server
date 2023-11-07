package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"
)

/*
FileIDGen returns an unique videoID and appends the fileExtension to it,
it takes the fileExtensionas parameter
*/
func FileIDGen(fileExtension string) string {
	length := 8
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	uid := base64.URLEncoding.EncodeToString(randomBytes)[:length]

	return uid + fileExtension
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
