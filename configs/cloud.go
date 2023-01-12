package configs

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type CloudSession struct {
	AWSSession   *session.Session
	GCPSession   *storage.Client
	AzureSession *azblob.SharedKeyCredential
}

var cloudSession *CloudSession

func InitCloud() {
	cloudSession = new(CloudSession)

	if EnvVar[AWS_ACCESS_KEY_ID] != "" {
		cloudSession.AWSSession = getAWSSession()
		log.Println("Initialised AWS")
	}

	if EnvVar[AZURE_ACCESS_KEY] != "" {
		cloudSession.AzureSession = getAzureSession()
		log.Println("Initialised Azure")
	}

	if EnvVar[GCP_PROJECT_ID] != "" {
		cloudSession.GCPSession = getGCPSession()
		log.Println("Initialised GCP")
	}
}

func GetCloudSession() *CloudSession {
	return cloudSession
}

func getAWSSession() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(EnvVar[AWS_S3_REGION]),
		Credentials: credentials.NewStaticCredentials(
			EnvVar[AWS_ACCESS_KEY_ID],
			EnvVar[AWS_SECRET_ACCESS_KEY],
			EnvVar[AWS_SESSION_TOKEN],
		),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})

	if err != nil {
		log.Println(err)
	}

	return s
}

func getAzureSession() *azblob.SharedKeyCredential {
	s, err := azblob.NewSharedKeyCredential(
		EnvVar[AZURE_ACCOUNT_NAME],
		EnvVar[AZURE_ACCESS_KEY],
	)

	if err != nil {
		log.Println(err)
	}

	return s
}

func getGCPSession() *storage.Client {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		log.Println(err)
	}

	return client
}
