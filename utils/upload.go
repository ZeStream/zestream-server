package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

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
	BucketName      string
	Prefix          string
	Region          string
	FolderLocalPath string
}

func (a AwsUploader) Upload(walker fileWalk) {

	//creating a new session
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(a.Region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to create a session")
	}

	bucket := a.BucketName
	log.Printf("bucket %s", bucket)

	prefix := a.Prefix

	uploader := s3manager.NewUploader(sess)
	for path := range walker {
		rel, err := filepath.Rel(a.FolderLocalPath, path)
		if err != nil {
			log.Fatalln("Unable to get relative path ", a.FolderLocalPath, path)
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

type GcpUploader struct {
	BucketName string
	ProjectId  string
	UploadPath string
}

func (g *GcpUploader) Upload(walker fileWalk) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("GCP_CREDENTIALS"))
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalln("Failed to create client ", err)
	}
	for path := range walker {
		filename := filepath.Base(path)
		fmt.Printf("Creating file /%v/%v\n", g.BucketName, filename)

		ctx := context.Background()

		wc := client.Bucket(g.BucketName).Object(g.UploadPath + filename).NewWriter(ctx)
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
	AzureEndpoint string
	ContainerName string
	AccountName   string
}

func (a AzureUploader) Upload(walker fileWalk) {

	for path := range walker {
		filename := filepath.Base(path)

		//create indiviual url for every blob
		u, _ := url.Parse(fmt.Sprint(a.AzureEndpoint, a.ContainerName, "/", filename))

		//create credential for
		credential, errC := azblob.NewSharedKeyCredential(a.AccountName, os.Getenv("AZURE_ACCESS_KEY"))
		if errC != nil {
			log.Fatalln("Failed to create credential")
		}
		blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

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
