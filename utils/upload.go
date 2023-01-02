package utils

import (
	"log"
	"os"
	"path/filepath"
	"zestream-server/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

type awsUploader struct {
	bucketName string
	prefix     string
}

func (a awsUploader) uploadToAWS(walker fileWalk, localpath string) {

	//creating a new session
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(constants.S_REGION),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to create a session")
	}

	bucket := a.bucketName
	log.Printf("bucket %s", bucket)
	if bucket == "" {
		log.Fatal("Error: AWS_S3_BUCKET env variable not set")
	}
	prefix := a.prefix
	if prefix == "" {
		log.Fatal("Error: AWS_S3_PREFIX env variable not set")
	}

	uploader := s3manager.NewUploader(sess)
	for path := range walker {
		rel, err := filepath.Rel(localpath, path)
		if err != nil {
			log.Fatalln("Unable to get relative path ", localpath, path)
		}

		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}

		defer file.Close()

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(filepath.Join(prefix, rel)),
			Body:   file,
		})

		if err != nil {
			log.Fatalln("Failed to upload", path, err)
		}
		log.Println("Uploaded", path, result.Location)
	}
}
