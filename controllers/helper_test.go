package controllers

import (
	"strings"
	"testing"
	"zestream-server/utils"

	"github.com/stretchr/testify/assert"
)

func TestVideoIDGen(t *testing.T) {

	// Call the Video ID Gen function
	videoIDExt := utils.VideoIDGen(".mp4")

	//Split the string into digits and ext
	videoID := strings.Split(videoIDExt, ".")

	//assert if the id is 8 digit long
	assert.Equal(t, 8, len(videoID[0]))
}
