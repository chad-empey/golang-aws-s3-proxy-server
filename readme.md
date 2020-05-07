# Golang AWS S3 Proxy Server

This is a simple golang s3 proxy server. Combine this with Google Cloud Run or any autoscaling infastructure and it's all you need to serve assets from an S3 bucket without needing the bucket to be public and other many benefits.

I've change the direction on this slightly and it currently does the folowing:

- Responds with a pre-signed url to upload directly to the bucket.
- Redirects to a pre-signed url to handle the serving of an asset in the bucket. The reason for this was to put more load on AWS and then you could possibly use CloudFront in front of the S3 bucket.
- Supports a delete request for an asset in the bucket.

# Usage

Dockerfile is included in the repository. All you need to do is set the following ENV variables for the AWS SDK.

- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_REGION
- AWS_BUCKET
