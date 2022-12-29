package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"zestream/constants"
)

func Fetch(url string, fileName string) error {
	newFileName, err := GetDownloadFilePathName(fileName)
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(constants.FILE_DOWNLOAD_ERROR + strconv.Itoa(response.StatusCode))
	}

	// create an empty file
	file, err := os.Create(newFileName)
	if err != nil {
		return err
	}

	defer file.Close()

	// write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
