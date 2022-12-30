package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"zestream-server/constants"
)

func Fetch(url string, fileName string) error {
	newFileName, err := getFileName(fileName)
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

func getFileName(fileName string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	newPath := path.Join(cwd, constants.DOWNLOAD_FILE_PATH_PREFIX)

	err = os.MkdirAll(newPath, os.ModePerm)

	if err != nil {
		return "", err
	}

	newPath = path.Join(newPath, fileName)

	return newPath, nil
}
