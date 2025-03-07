package main

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := dynamo.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %+v", err)
	}

	err = db.CreateTestTable()
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "ResourceInUseException" {
			log.Info("Error is okay. Table already exists 👍")
		} else {
			log.Fatalf("failed to create table: %+v", err)
		}
	}

	s3, err := s3.NewS3()
	if err != nil {
		log.Fatalf("failed to connect to s3: %+v", err)
	}

	err = s3.CreateTestBucket()
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "BucketAlreadyOwnedByYou" {
			log.Info("Error is okay. Bucket already exists 👍")
			return
		}
		log.Fatalf("failed to create bucket: %+v", err)
	}

	log.Info("done")
}
