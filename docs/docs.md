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