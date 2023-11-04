package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileIDGen(t *testing.T) {

	// Call the Video ID Gen function
	videoIDExt := FileIDGen(".mp4")

	//Split the string into digits and ext
	videoID := strings.Split(videoIDExt, ".")

	//assert if the id is 8 digit long
	assert.Equal(t, 8, len(videoID[0]))
}
