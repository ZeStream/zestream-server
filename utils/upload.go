package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"zestream-server/constants"

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
			log.Fatalln("Walk failed: ", err)
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
	Prefix string
	//common session to be used by every upload
	Session *session.Session
}

func (a AwsUploader) Upload(walker fileWalk) {
	bucket := constants.S3_BUCKET_NAME
	if bucket == "" {
		log.Fatalln("AWS Bucketname not available")
	}
	log.Printf("bucket %s", bucket)

	prefix := a.Prefix

	uploader := s3manager.NewUploader(a.Session)
	for path := range walker {
		filename := filepath.Base(path)

		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}

		defer file.Close()

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(filepath.Join(prefix, filename)),
			Body:   file,
		})

		if err != nil {
			log.Fatalln("Failed to upload", path, err)
		}
		log.Println("Uploaded", path, result.Location)
	}
}

type GcpUploader struct {
	UploadPath string
	//common azure storage client to be used for every upload
	Client *storage.Client
}

func (g *GcpUploader) Upload(walker fileWalk) {
	bucketName := constants.GCP_BUCKET_NAME
	if bucketName == "" {
		log.Fatalln("GCP Bucketname not available")
	}
	for path := range walker {
		filename := filepath.Base(path)
		fmt.Printf("Creating file /%v/%v\n", bucketName, filename)

		ctx := context.Background()

		wc := g.Client.Bucket(bucketName).Object(g.UploadPath + filename).NewWriter(ctx)
		blob, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
		}
		defer blob.Close()
		if _, err := io.Copy(wc, blob); err != nil {
			log.Fatalln("Failed to upload", path, err)
		}

		if err := wc.Close(); err != nil {
			log.Fatalln("unable to close the bucket", err)
		}
	}

}

type AzureUploader struct {
	ContainerName string

	//common for every upload process
	AzureCredential *azblob.SharedKeyCredential
}

func (a AzureUploader) Upload(walker fileWalk) {
	accountName := constants.AZURE_ACCOUNT_NAME
	azureEndpoint := constants.AZURE_ENDPOINT
	if accountName == "" {
		log.Fatalln("azure account name not available")
	}
	if azureEndpoint == "" {
		log.Fatalf("azure endpoint not available")
	}
	for path := range walker {
		filename := filepath.Base(path)

		//create indiviual url for every blob
		u, _ := url.Parse(fmt.Sprint(azureEndpoint, a.ContainerName, "/", filename))
		blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(a.AzureCredential, azblob.PipelineOptions{}))

		ctx := context.Background()

		// Upload to data to blob storage
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed to open file ", path)
			continue
		}
		defer file.Close()
		_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobUrl, azblob.UploadToBlockBlobOptions{})
		if err != nil {
			log.Fatalln("Failure to upload to azure container:")
		} else {
			log.Printf("successfully uploaded %s ", path)
		}
	}

}
