package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"zestream-server/configs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"cloud.google.com/go/storage"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// channel to extract files from the folder
type fileWalk chan string

type Uploader interface {
	Upload(walker fileWalk)
}

func UploadToCloudStorage(uploader Uploader, path string) {
	walker := make(fileWalk)

	go func() {
		//get files to upload via the channel
		if err := filepath.Walk(path, walker.WalkFunc); err != nil {
			log.Println("Walk failed: ", err)
		}

		close(walker)
	}()

	uploader.Upload(walker)
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

type AwsUploader struct {
	Prefix  string
	VideoId string
	Session *session.Session
}

func (a AwsUploader) Upload(walker fileWalk) {
	bucket := configs.EnvVar[configs.AWS_S3_BUCKET_NAME]

	if bucket == "" {
		log.Println("AWS Bucketname not available")
	}

	uploader := s3manager.NewUploader(a.Session)
	for path := range walker {
		filename := filepath.Base(path)

		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(filepath.Join(a.Prefix, a.VideoId, filename)),
			Body:   file,
		})

		if err != nil {
			file.Close()
			log.Println("Failed to upload", path, err)
		}
		log.Println("Uploaded", path, result.Location)

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file")
		}
	}
}

type GcpUploader struct {
	UploadPath string
	VideoId    string
	Client     *storage.Client
}

func (g *GcpUploader) Upload(walker fileWalk) {
	bucketName := configs.EnvVar[configs.GCP_BUCKET_NAME]
	if bucketName == "" {
		log.Println("GCP Bucketname not available")
	}

	for path := range walker {
		filename := filepath.Base(path)
		fmt.Printf("Creating file /%v/%v\n", bucketName, filename)

		ctx := context.Background()

		wc := g.Client.Bucket(bucketName).Object(filepath.Join(g.UploadPath, g.VideoId, filename)).NewWriter(ctx)
		blob, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
		}

		if _, err := io.Copy(wc, blob); err != nil {
			//close the blob before error log
			blob.Close()
			log.Println("Failed to upload", path, err)
		}

		if err := wc.Close(); err != nil {
			//close the file before error log
			blob.Close()
			log.Println("unable to close the bucket", err)
		} else {
			log.Println("successfully uploaded ", path)
		}

		if err := blob.Close(); err != nil {
			log.Println("unable to close the file")
		}
	}

}

type AzureUploader struct {
	ContainerName string
	VideoId       string

	Credential *azblob.SharedKeyCredential
}

func (a AzureUploader) Upload(walker fileWalk) {
	azureEndpoint := configs.EnvVar[configs.AZURE_ENDPOINT]

	log.Println(azureEndpoint)

	if azureEndpoint == "" {
		log.Println("Azure endpoint not available")
	}

	for path := range walker {
		filename := filepath.Base(path)

		//create indiviual url for every blob
		u, _ := url.Parse(azureEndpoint)
		u = u.JoinPath(a.ContainerName, a.VideoId, filename)

		blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(a.Credential, azblob.PipelineOptions{}))

		ctx := context.Background()

		// Upload to data to blob storage
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed to open file ", path)
			continue
		}

		_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobUrl, azblob.UploadToBlockBlobOptions{})
		if err != nil {
			//close the file before error log
			defer file.Close()
			log.Println("Failure to upload to azure container:", err)
		} else {
			log.Printf("successfully uploaded %s ", path)
		}

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file ", path)
		}
	}

}
