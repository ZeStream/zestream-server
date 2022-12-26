package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"zestream/constants"
)

func GetDownloadFilePathName(fileName string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	pathName := path.Join(cwd, constants.DOWNLOAD_FILE_PATH_PREFIX)

	err = createDirPath(pathName)

	if err != nil {
		return "", err
	}

	newPath := path.Join(pathName, fileName)

	return newPath, nil
}

func GetOutputFilePathName(fileName string, postfix string) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	// TODO: replace filename with id
	fileName = strings.Replace(fileName, ".", "_", 1)

	pathName := path.Join(cwd, constants.OUTPUT_FILE_PATH_PREFIX, fileName)

	err = createDirPath(pathName)

	fmt.Println(pathName)

	if err != nil {
		return "", err
	}

	newPath := path.Join(pathName, postfix)

	return newPath, nil
}

func createDirPath(pathName string) error {

	err := os.MkdirAll(pathName, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}
