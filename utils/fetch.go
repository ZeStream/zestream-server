package utils

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"zestream-server/constants"
)

/*
Fetch downlods a file to downloads folder from the given url,
and names it as given fileName
*/
func Fetch(url string, fileName string) error {
	log.Println("Downloading: ", url)

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

	log.Println("Download Complete: ", url)

	return nil
}
