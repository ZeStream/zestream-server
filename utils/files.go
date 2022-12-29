package utils

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

func RemoveExtensionFromFile(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func IsFileValid(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func WrapStringInQuotes(str string) string {
	var buff bytes.Buffer

	buff.WriteString(str)
	buff.WriteString(" ")

	return buff.String()
}

func StringToArgsGenerator(args map[string]string) string {
	var argsStr bytes.Buffer

	for k, v := range args {
		argsStr.WriteString("-")
		argsStr.WriteString(k + " ")
		argsStr.WriteString(v)

		if v != "" {
			argsStr.WriteString(" ")
		}
	}

	return argsStr.String()
}

func DeleteFiles(filePaths string) {
	_, err := exec.Command("rm", strings.Split(filePaths, " ")...).Output()

	if err != nil {
		log.Println(err)
	}
}
