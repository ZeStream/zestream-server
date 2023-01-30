## Configuring AWS

1. Set up AWS account using [this](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/#:~:text=Sign%20up%20using%20your%20email,Create%20a%20new%20AWS%20account.)

2. You will need `access key ID `, `secret access key`, `session token`, `S3 region` and `S3 bucket name`
   - [Get Access Keys](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-about)
   - [Get Session Token](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html)
   - [Creating S3 Bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/create-bucket-overview.html)

3. Set the values in the `.env` file
   | Env Variable | Value |  
   |--------------|-------|
   |AWS_ACCESS_KEY_ID | access key ID |  
   | AWS_SECRET_ACCESS_KEY | secret access key |  
   | AWS_SESSION_TOKEN | session token |
   | AWS_S3_BUCKET_NAME | S3 bucket name |
   | AWS_S3_REGION | S3 region |

## Configuring GCP

1. Set up GCP account using [this](https://cloud.google.com/apigee/docs/hybrid/v1.2/precog-gcpaccount)

2. You will need `project ID` and `GCP bucket name`
   - [Get Project ID](https://support.google.com/googleapi/answer/7014113?hl=en)
   - [Creating GCP Bucket](https://cloud.google.com/storage/docs/creating-buckets)

3. Set the values in the `.env` file
   | Env Variable | Value |  
   |--------------|-------|
   |GCP_PROJECT_ID| project ID|
   |GCP_BUCKET_NAME | bucket name |

## Configuring Azure

1. Create a storage account using [this](https://azure.microsoft.com/en-in/free/)

2. You will need `account name` , `azure endpoint` and `access key`

   - [Create a container](https://learn.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-portal#create-a-container)
   - [Get Access Key](https://learn.microsoft.com/en-us/azure/storage/common/storage-account-keys-manage?tabs=azure-portal)
   - If your storage account is named 'mystorageaccount', then the default endpoint for Blob Storage is:

   ```

       http://mystorageaccount.blob.core.windows.net

   ```

3. Set the values in the `.env` file
   | Env Variable | Value |  
    |--------------|-------|
   |AZURE_ACCESS_KEY| access key|
   |AZURE_ENDPOINT | azure endpoint |
   |AZURE_ACCOUNT_NAME | name of the storage account |

