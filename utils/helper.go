package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"log"
	"sort"
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

// sortString sorts the given string
func SortString(input string) string {
	runes := []rune(input)

	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return string(runes)
}

// GetQueryHash returns the md5 hash of the query string
func GetQueryHash(query string) string {
	sortedQuery := SortString(query)

	hasher := md5.New()
	hasher.Write([]byte(sortedQuery))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}
