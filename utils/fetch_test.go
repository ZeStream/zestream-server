package utils_test

import (
	"testing"
	"zestream-server/utils"
)

func TestFetch(t *testing.T) {
	const fileURL = "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/WeAreGoingOnBullrun.mp4"
	const fileName = "WeAreGoingOnBullrun.mp4"

	// Download the file locally
	err := utils.Fetch(fileURL, fileName)

	if err != nil {
		t.Fatal(err)
	}

	pathName, err := utils.GetDownloadFilePathName(fileName)

	if err != nil {
		t.Fatal(err)
	}

	isFileDownloaded := utils.IsFileValid(pathName)

	if isFileDownloaded == false {
		t.Error("want true, got", isFileDownloaded)
	}

	t.Error(isFileDownloaded, pathName)

	utils.DeleteFiles(pathName)
}
