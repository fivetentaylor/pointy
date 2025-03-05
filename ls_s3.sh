#!/bin/bash

# Set default bucket name from environment variable or default
# s3://reviso-documents
BUCKET_NAME="${S3_BUCKET:-reviso-documents}"

# Set prefix from environment variable if provided
PREFIX="${1:-}"

aws s3 ls "s3://${BUCKET_NAME}/${PREFIX}" --recursive |\
   awk '{print $4}' |\
   parallel "echo 's3://${BUCKET_NAME}/{}'"
