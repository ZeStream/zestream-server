package utils

import (
	"log"
	"os"
	"path/filepath"
)

// channel to extract files from the folder
type fileWalk chan string

func UploadToCloudStorage(path string, cloudPlatform string) {
	walker := make(fileWalk)
	go func() {
		//get files to upload via the channel
		if err := filepath.Walk(path, walker.WalkFunc); err != nil {
			log.Fatalln("Walk failed: ", err)
		}

		close(walker)

	}()
}

func (f fileWalk) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		f <- path
	}

	return nil
}
