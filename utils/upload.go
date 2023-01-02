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

type gcpUploader struct {
	bucketName string
	projectId  string
	client     *storage.Client
	uploadPath string
}

func (g *gcpUploader) uploadtToGcp(walker fileWalk) {

	for path := range walker {
		filename := filepath.Base(path)
		fmt.Printf("Creating file /%v/%v\n", g.bucketName, filename)

		ctx := context.Background()

		wc := g.client.Bucket(g.bucketName).Object(g.uploadPath + filename).NewWriter(ctx)
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

type azureUploader struct {
	azureEndpoint string
	containerName string
	accountName   string
}

func (a azureUploader) uploadToAzure(walker fileWalk) {

	for path := range walker {
		filename := filepath.Base(path)

		//create indiviual url for every blob
		u, _ := url.Parse(fmt.Sprint(a.azureEndpoint, a.containerName, "/", filename))

		//create credential for
		credential, errC := azblob.NewSharedKeyCredential(a.accountName, os.Getenv("AZURE_ACCESS_KEY"))
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
