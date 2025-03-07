package s3

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

type S3 struct {
	Bucket       string
	ImagesBucket string
	Client       *s3.S3
	AppHost      string
}

func NewS3() (*S3, error) {
	awsConfig := aws.Config{}

	s3_bucket, exists := os.LookupEnv("AWS_S3_BUCKET")
	if !exists {
		return nil, fmt.Errorf("AWS_S3_BUCKET not set in environment variables")
	}

	s3_images_bucket, exists := os.LookupEnv("AWS_S3_IMAGES_BUCKET")
	if !exists {
		return nil, fmt.Errorf("AWS_S3_IMAGES_BUCKET not set in environment variables")
	}

	s3_url, exists := os.LookupEnv("AWS_S3_URL")
	if exists {
		// We're in test/dev
		log.Printf("S3 URL: %s", s3_url)
		awsConfig.Endpoint = aws.String(s3_url)
		// If we're here this means we're in test and should
		// use path style because it's hard to setup in test
		// subdomain style is what prod aws requires now
		awsConfig.S3ForcePathStyle = aws.Bool(true)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		awsConfig.HTTPClient = &http.Client{Transport: tr}
	}

	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	return &S3{
		Bucket:       s3_bucket,
		ImagesBucket: s3_images_bucket,
		Client:       s3.New(sess),
		AppHost:      os.Getenv("APP_HOST"),
	}, nil

}

func (s *S3) _createTestBucket(name string) error {
	_, err := s.Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "BucketAlreadyOwnedByYou" {
			return nil
		}
		log.Errorf("failed to create bucket: %s", err)
		return err
	}

	err = s.Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		log.Errorf("failed waiting on bucket exists: %s", err)
		return err
	}

	return nil

}

func (s *S3) CreateTestBucket() error {
	err := s._createTestBucket(s.Bucket)
	if err != nil {
		return err
	}

	err = s._createTestBucket(s.ImagesBucket)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) PutObject(bucket, key, contentType string, body []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(body),
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return stackerr.New(err)
	}

	return nil
}

func (s *S3) PutObjectWithSeeker(bucket, key, contentType string, body io.ReadSeeker) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Body:        body,
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Errorf("Failed to put object: %s", err)
		return err
	}

	return nil
}

func (s *S3) PutCachedObject(bucket, key, contentType string, body []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Body:         bytes.NewReader(body),
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("max-age=86400"),
	})
	if err != nil {
		log.Errorf("Failed to put object: %s", err)
		return err
	}

	return nil
}

func (s *S3) DeletePrefix(bucket, prefix string) error {
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	log.Debugf("Deleting prefix %s: %#v", prefix, listObjectsInput)

	err := s.Client.ListObjectsV2Pages(listObjectsInput,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			if len(page.Contents) == 0 {
				fmt.Println("No objects found")
				return false
			}

			var objectsToDelete []*s3.ObjectIdentifier
			for _, object := range page.Contents {
				objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
					Key: object.Key,
				})
			}

			log.Debugf("Deleting %#v objects", objectsToDelete)

			deleteObjectsInput := &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3.Delete{
					Objects: objectsToDelete,
					Quiet:   aws.Bool(true),
				},
			}

			deleteObjectsOutput, err := s.Client.DeleteObjects(deleteObjectsInput)
			if err != nil {
				log.Fatalf("Failed to delete objects: %v", err)
			}

			log.Debugf("Deleted %d objects", len(deleteObjectsOutput.Deleted))

			return true
		})
	if err != nil {
		log.Fatalf("Failed to list objects: %v", err)
		return err
	}

	return nil
}

func (s *S3) ObjectExists(bucket, key string) (bool, error) {
	_, err := s.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" {
				return false, nil
			}
			return false, fmt.Errorf("AWS ERROR CODE: %s, ERROR: %s", aerr.Code(), aerr.Message())
		}
		return false, err
	}

	return true, nil
}

func (s *S3) GetObject(bucket, key string) ([]byte, error) {
	obj, err := s.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Errorf("Failed to get object: %s", err)
		return nil, err
	}
	defer obj.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		log.Errorf("Failed to read into buffer: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *S3) GetAllObjects(bucket, prefix string) ([][]byte, error) {
	objects := [][]byte{}

	err := s.StreamAllObjects(bucket, prefix, func(b []byte) error {
		objects = append(objects, b)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (s *S3) StreamAllObjects(bucket, prefix string, callback func([]byte) error) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	for {
		result, err := s.Client.ListObjectsV2(input)
		if err != nil {
			log.Errorf("Failed to list objects: %s", err)
			return err
		}

		// Process the current page of objects
		for _, object := range result.Contents {
			obj, err := s.Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    object.Key,
			})
			if err != nil {
				return err
			}
			defer obj.Body.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(obj.Body)
			if err != nil {
				return err
			}

			err = callback(buf.Bytes())
			if err != nil {
				return err
			}
		}

		// If there are more objects, set the ContinuationToken for the next page
		if result.IsTruncated != nil && *result.IsTruncated {
			input.ContinuationToken = result.NextContinuationToken
		} else {
			break
		}
	}

	return nil

}

func (s *S3) List(bucket, prefix string, pageIx, pageSize int) ([]string, error) {
	if pageSize < 1 {
		pageSize = 1000
	}

	var out []string
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(int64(pageSize)),
	}

	currentPage := 0

	err := s.Client.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			currentPage++ // Increment the current page counter

			// Check if the current page is the one we want
			if pageIx < 0 || currentPage == pageIx {
				for _, obj := range page.Contents {
					out = append(out, *obj.Key)
				}

				if pageIx > 0 {
					return false
				}
			}

			// If it's not the desired page, check if this is the last page; if not, continue
			return !lastPage
		})

	if err != nil {
		log.Errorf("failed to list objects: %s", err)
		return nil, err
	}

	return out, nil
}

func (s *S3) DeleteAll(bucket, prefix string) error {
	objects, err := s.List(bucket, prefix, 1, 1000)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(obj),
		})
		if err != nil {
			log.Errorf("failed to delete object: %s", err)
			return err
		}
	}

	return nil
}

func (s *S3) GetPresignedUrl(bucket, key string, expire time.Duration) (string, error) {
	req, _ := s.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expire)
	if err != nil {
		fmt.Println("Failed to sign request", err)
		return "", err
	}

	return urlStr, nil
}

func (s *S3) ListFirst(bucket, prefix string) (string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(1),
	}

	result, err := s.Client.ListObjectsV2(input)
	if err != nil {
		return "", fmt.Errorf("failed to list objects: %s", err)
	}

	if len(result.Contents) == 0 {
		return "", fmt.Errorf("prefix not found: %s", prefix)
	}

	return *result.Contents[0].Key, nil
}

func (s *S3) ListFirstTwo(bucket, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(2),
	}

	result, err := s.Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %s", err)
	}

	if len(result.Contents) < 2 {
		return nil, fmt.Errorf("prefix not found: %s", prefix)
	}

	out := []string{*result.Contents[0].Key, *result.Contents[1].Key}
	return out, nil
}

func (s *S3) Exists(bucket, key string) (bool, error) {
	_, err := s.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}
