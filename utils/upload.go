package utils

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"zestream-server/constants"
	"zestream-server/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"cloud.google.com/go/storage"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// channel to extract files from the folder
type fileWalk chan string

type Uploader interface {
	Upload(ctx context.Context, walker fileWalk)
}

func UploadToCloudStorage(ctx context.Context, uploader Uploader, path string) {
	walker := make(fileWalk)
	go func() {
		//get files to upload via the channel
		if err := filepath.Walk(path, walker.WalkFunc); err != nil {
			logger.Error(ctx, "Walk Failed", logger.Z{
				"error": err.Error(),
				"path":  path,
			})
		}

		close(walker)

	}()

	uploader.Upload(ctx, walker)
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

func (a AwsUploader) Upload(ctx context.Context, walker fileWalk) {
	bucket := constants.S3_BUCKET_NAME
	if bucket == "" {
		logger.Info(ctx, "AWS Bucketname not available", logger.Z{})
	}

	logger.Info(ctx, "AWS bucketname", logger.Z{
		"bucket_name": bucket,
	})

	prefix := a.Prefix

	uploader := s3manager.NewUploader(a.Session)
	for path := range walker {
		filename := filepath.Base(path)

		file, err := os.Open(path)
		if err != nil {
			logger.Error(ctx, "Failed opening file", logger.Z{
				"error": err.Error(),
				"path":  path,
			})

			continue
		}

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(filepath.Join(prefix, filename)),
			Body:   file,
		})

		if err != nil {
			//close the file before error log
			file.Close()
			logger.Error(ctx, "Failed to upload", logger.Z{
				"error":     err.Error(),
				"prefix":    prefix,
				"file_name": filename,
				"bucket":    bucket,
			})
		}

		logger.Info(ctx, "Video Uploaded successfully", logger.Z{
			"path":   path,
			"result": result,
		})

		if err := file.Close(); err != nil {
			logger.Error(ctx, "Unable to close file", logger.Z{
				"error": err.Error(),
				"path":  path,
			})
		}
	}
}

type GcpUploader struct {
	UploadPath string
	//common azure storage client to be used for every upload
	Client *storage.Client
}

func (g *GcpUploader) Upload(ctx context.Context, walker fileWalk) {
	bucketName := constants.GCP_BUCKET_NAME
	if bucketName == "" {
		logger.Error(ctx, "GCP bucketname not available", logger.Z{})
	}

	for path := range walker {
		filename := filepath.Base(path)
		logger.Info(ctx, "Creating file", logger.Z{
			"bucket_name": bucketName,
			"file_name":   filename,
		})

		wc := g.Client.Bucket(bucketName).Object(g.UploadPath + filename).NewWriter(ctx)
		blob, err := os.Open(path)
		if err != nil {
			logger.Error(ctx, "Failed opening file", logger.Z{
				"error": err.Error(),
				"path":  path,
			})

		}

		if _, err := io.Copy(wc, blob); err != nil {
			//close the blob before error log
			blob.Close()
			logger.Error(ctx, "Failed to upload", logger.Z{
				"error":       err.Error(),
				"file_name":   filename,
				"upload_path": g.UploadPath,
			})
		}

		if err := wc.Close(); err != nil {
			//close the file before error log
			blob.Close()
			logger.Error(ctx, "Unable to close bucket", logger.Z{
				"error":       err.Error(),
				"bucker_name": bucketName,
			})
		} else {
			logger.Info(ctx, "successfully uploaded", logger.Z{
				"path": path,
			})
		}

		if err := blob.Close(); err != nil {
			logger.Error(ctx, "unable to close file", logger.Z{
				"error":     err.Error(),
				"file_name": filename,
			})
		}
	}
}

type AzureUploader struct {
	ContainerName string

	//common for every upload process
	AzureCredential *azblob.SharedKeyCredential
}

func (a AzureUploader) Upload(ctx context.Context, walker fileWalk) {
	accountName := constants.AZURE_ACCOUNT_NAME
	azureEndpoint := constants.AZURE_ENDPOINT
	if accountName == "" {
		logger.Error(ctx, "azure acount name not available", logger.Z{})
	}
	if azureEndpoint == "" {
		logger.Error(ctx, "azure endpoint not available", logger.Z{})
	}

	for path := range walker {
		filename := filepath.Base(path)

		//create indiviual url for every blob
		u, _ := url.Parse(fmt.Sprint(azureEndpoint, a.ContainerName, "/", filename))
		blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(a.AzureCredential, azblob.PipelineOptions{}))

		// Upload to data to blob storage
		file, err := os.Open(path)
		if err != nil {
			logger.Info(ctx, "Failed to open file", logger.Z{
				"error": err.Error(),
				"path":  path,
			})

			continue
		}

		_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobUrl, azblob.UploadToBlockBlobOptions{})
		if err != nil {
			//close the file before error log
			file.Close()
			logger.Error(ctx, "Failure in uploading to azure container", logger.Z{
				"error": err.Error(),
				"path":  path,
			})
		} else {
			logger.Info(ctx, "successfully uploaded", logger.Z{
				"path": path,
			})
		}

		if err := file.Close(); err != nil {
			logger.Error(ctx, "Unable to close the file", logger.Z{
				"error": err.Error(),
				"path":  path,
			})
		}
	}

}
