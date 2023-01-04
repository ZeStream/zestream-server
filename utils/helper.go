package utils

import (
	"crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
	"strconv"
)

func VideoIDGen(fileExtension string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return err.Error()
	}

	var i int64 = int64(binary.LittleEndian.Uint64(b[:]))
	math_rand.Seed(i)

	// Generate a 8 digit random number
	randomNumber := math_rand.Intn(100000000) + 10000000

	// VideoID = (8-digit random number) + (file Name)
	return strconv.Itoa(randomNumber) + fileExtension
}
