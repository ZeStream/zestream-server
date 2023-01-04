package utils

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"math/rand"
	"strconv"
)

func VideoIDGen(fileExtension string) string {
	// Seed the random number generator
	hash := crc64.New(crc64.MakeTable(crc64.ECMA))

	// Write the string to the hash
	_, err := hash.Write([]byte("zestream"))
	if err != nil {
		fmt.Println(err)
	}

	// Calculate the hash
	sum := hash.Sum(nil)

	// Convert the hash to an int64 value
	var i int64
	i = int64(binary.LittleEndian.Uint64(sum))
	rand.Seed(i)

	// Generate a 8 digit random number
	randomNumber := rand.Intn(100000000) + 10000000

	// VideoID = (8-digit random number) + (file Name)
	return strconv.Itoa(randomNumber) + fileExtension
}
