package configs

import (
	"context"
	"encoding/base64"
	"log"

	"cloud.google.com/go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"google.golang.org/api/option"
)

type CloudSession struct {
	AWS   *session.Session
	GCP   *storage.Client
	Azure *azblob.SharedKeyCredential
}

var cloudSession *CloudSession

func InitCloud() {
	cloudSession = new(CloudSession)

	if EnvVar[AWS_ACCESS_KEY_ID] != "" {
		cloudSession.AWS = getAWSSession()
		log.Println("Initialised AWS")
		return
	}

	if EnvVar[AZURE_ACCESS_KEY] != "" {
		cloudSession.Azure = getAzureCreds()
		log.Println("Initialised Azure")
		return
	}

	if EnvVar[GCP_PROJECT_ID] != "" {
		cloudSession.GCP = GetGCPClient(false)
		log.Println("Initialised GCP")
		return
	}

	log.Fatalln("Cloud Initialization Failed")
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

func getAzureCreds() *azblob.SharedKeyCredential {
	s, err := azblob.NewSharedKeyCredential(
		EnvVar[AZURE_ACCOUNT_NAME],
		EnvVar[AZURE_ACCESS_KEY],
	)

	if err != nil {
		log.Println(err)
	}

	return s
}

func GetGCPClient(isServiceUser bool) *storage.Client {
	var credsInBase64 string

	if isServiceUser {
		credsInBase64 = EnvVar[GCP_SERVICE_USER_CREDS_JSON_BASE64]
	} else {
		credsInBase64 = EnvVar[GCP_CREDS_JSON_BASE64]
	}

	credsDecoded, err := base64.RawURLEncoding.DecodeString(credsInBase64)
	if err != nil {
		log.Println(err)
	}

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(credsDecoded))
	if err != nil {
		log.Println(err)
	}

	return client
}
