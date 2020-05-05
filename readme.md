# Golang AWS S3 Proxy Server

This is a powerful golang s3 proxy server. Combine this with Google Cloud Run or any autoscaling infastructure and it's all you need to serve assets from an S3 bucket without needing the bucket to be public and other many benefits.

# Usage

Dockerfile is included in the repository. All you need to do is set the following ENV variables for the AWS SDK.

- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_REGION
- AWS_BUCKET

# Useful commands for Cloud Run

If you're using Google Cloud Run with this and you have installed the gcloud cli, then the following commands are easy to build and deploy.

- gcloud builds submit --tag gcr.io/YOUR_INFO/IMAGE_NAME
- gcloud run deploy --image gcr.io/YOUR_INFO/IMAGE_NAME

You could also easily automate the build/deploy process.
