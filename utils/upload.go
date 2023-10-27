package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
	"zestream-server/configs"
	"zestream-server/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"cloud.google.com/go/storage"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// channel to extract files from the folder
type fileWalk chan string

type Uploader interface {
	Upload(walker fileWalk)
}

func GetUploader(containerName string, videoId string) Uploader {

	cloudSession := configs.GetCloudSession()
	isAWS, isGCP, isAzure := GetWhichCloudIsEnabled()

	if isAWS {
		return AwsUploader{
			ContainerName: containerName,
			VideoId:       videoId,
			Session:       cloudSession.AWS,
		}
	}

	if isGCP {
		return &GcpUploader{
			ContainerName: containerName,
			VideoId:       videoId,
			Client:        cloudSession.GCP,
		}
	}

	if isAzure {
		return AzureUploader{
			ContainerName: containerName,
			VideoId:       videoId,
			Credential:    cloudSession.Azure,
		}
	}

	return nil
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
	ContainerName string
	VideoId       string
	Session       *session.Session
}

func (a AwsUploader) Upload(walker fileWalk) {
	bucket := configs.EnvVar[configs.AWS_S3_BUCKET_NAME]

	if bucket == "" {
		log.Println("AWS Bucketname not available")
	}

	uploader := s3manager.NewUploader(a.Session)

	for pathName := range walker {
		filename := filepath.Base(pathName)

		file, err := os.Open(pathName)
		if err != nil {
			log.Println("Failed opening file", pathName, err)
			continue
		}

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(path.Join(a.ContainerName, a.VideoId, filename)),
			Body:   file,
		})

		if err != nil {
			file.Close()
			log.Println("Failed to upload", pathName, err)
		}
		log.Println("Uploaded", pathName, result.Location)

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file")
		}
	}
}

type GcpUploader struct {
	ContainerName string
	VideoId       string
	Client        *storage.Client
}

func (g *GcpUploader) Upload(walker fileWalk) {
	bucketName := configs.EnvVar[configs.GCP_BUCKET_NAME]
	if bucketName == "" {
		log.Println("GCP Bucketname not available")
	}

	ctx := context.Background()

	bucket := g.Client.Bucket(bucketName)

	for pathName := range walker {
		filename := filepath.Base(pathName)
		fmt.Printf("Creating file /%v/%v\n", bucketName, filename)

		wc := bucket.Object(path.Join(g.ContainerName, g.VideoId, filename)).NewWriter(ctx)

		blob, err := os.Open(pathName)
		if err != nil {
			log.Println("Failed opening file", pathName, err)
		}

		if _, err := io.Copy(wc, blob); err != nil {
			//close the blob before error log
			blob.Close()
			log.Println("Failed to upload", pathName, err)
		}

		if err := wc.Close(); err != nil {
			//close the file before error log
			blob.Close()
			log.Println("unable to close the bucket", err)
		} else {
			log.Println("successfully uploaded ", pathName)
		}

		if err := blob.Close(); err != nil {
			log.Println("unable to close the file")
		}
	}

}

type AzureUploader struct {
	ContainerName string
	VideoId       string
	Credential    *azblob.SharedKeyCredential
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
			log.Println("Failure to upload to azure container:", err)
			file.Close()
		} else {
			log.Printf("successfully uploaded %s ", path)
		}

		if err := file.Close(); err != nil {
			log.Println("Unable to close the file ", path)
		}
	}

}

func GetWhichCloudIsEnabled() (bool, bool, bool) {
	if configs.EnvVar[configs.AWS_ACCESS_KEY_ID] != "" {
		return true, false, false
	}

	if configs.EnvVar[configs.GCP_BUCKET_NAME] != "" {
		return false, true, false
	}

	if configs.EnvVar[configs.AZURE_ACCOUNT_NAME] != "" {
		return false, false, true
	}

	return false, false, false
}

func GetSignedURL(videoId string, fileType string, basePath string) string {
	isAWS, isGCP, isAzure := GetWhichCloudIsEnabled()

	filePath, err := GetCloudStoragePath(basePath, videoId, fileType)
	LogErr(err)

	if isAWS {
		return generateAWSSignedURL(filePath)
	}

	if isGCP {
		return generateGCPSignedURL(filePath)
	}

	if isAzure {
		return generateAzureSignedURL(filePath)
	}

	return ""
}

func generateGCPSignedURL(filePath string) string {
	client := configs.GetGCPClient(true)
	bucket := client.Bucket(configs.EnvVar[configs.GCP_BUCKET_NAME])

	expirationTime := time.Now().Add(constants.PRESIGNED_URL_EXPIRATION)

	url, err := bucket.SignedURL(filePath, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: expirationTime,
	})
	if err != nil {
		log.Println(err)
	}

	return url
}

func generateAWSSignedURL(filePath string) string {
	session := configs.GetCloudSession()

	client := s3.New(session.AWS)

	req, err := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(configs.EnvVar[configs.AWS_S3_BUCKET_NAME]),
		Key:    aws.String(filePath),
	})

	if err != nil {
		log.Println(err)
	}

	url, er := req.Presign(constants.PRESIGNED_URL_EXPIRATION)
	if er != nil {
		log.Println(err)
	}

	return url
}

func generateBlobSAS(blobURL azblob.BlobURL, permissions string) (string, error) {
	session := configs.GetCloudSession()
	client := session.Azure

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:    azblob.SASProtocolHTTPS,
		ExpiryTime:  time.Now().Add(constants.PRESIGNED_URL_EXPIRATION),
		Permissions: permissions,
		Version:     time.Now().Format(time.DateOnly),
	}.NewSASQueryParameters(client)

	if err != nil {
		return "", err
	}

	sasURL := blobURL.URL()
	sasURL.RawQuery = sasQueryParams.Encode()
	log.Println(sasQueryParams.Encode())
	return sasURL.String(), nil
}

func generateAzureSignedURL(filePath string) string {
	accountName := configs.EnvVar[configs.AZURE_ACCOUNT_NAME]
	accountKey := configs.EnvVar[configs.AZURE_ACCESS_KEY]
	containerName := constants.CloudContainerNames[constants.Images]

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println("Failed to create shared key credential:", err)
		return ""
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	blobURL := azblob.NewContainerURL(*URL, pipeline).NewBlobURL(filePath)

	permissions := "rw"

	sasURL, err := generateBlobSAS(blobURL, permissions)
	if err != nil {
		fmt.Println("Failed to generate pre-signed URL:", err)
		return ""
	}

	return sasURL
}

func GetCloudStoragePath(basePath, fileName string, _ string) (string, error) {
	url, err := url.JoinPath(basePath, fileName)
	LogErr(err)

	return url, err
}
