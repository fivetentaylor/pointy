package document

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

const screenshotWidth = int64(640)
const screenshotHeight = int64(600)
const screenshotUrlExpireDuration = 30 * time.Minute

func SeqToScreenshotS3Path(docID string, seq int64, theme string) string {
	return path.Join(
		constants.S3Prefix,
		docID,
		"screenshots",
		fmt.Sprintf("%s-%s.png", utils.InvertSeq(seq), theme),
	)
}

func GetEmailScreenshotsURL(ctx context.Context, docID string) (string, error) {
	s3 := env.S3(ctx)
	screenPath := path.Join(constants.S3Prefix, docID, "screenshots")
	screenKeys, err := s3.ListFirstTwo(s3.Bucket, screenPath)
	if err != nil {
		return "", fmt.Errorf("error getting last screenshot: %s", err)
	}

	lightUrl, err := s3.GetPresignedUrl(
		s3.Bucket,
		screenKeys[1],
		screenshotUrlExpireDuration,
	)
	if err != nil {
		return "", err
	}

	return lightUrl, nil
}

func GetLastScreenshotsURL(ctx context.Context, docID string) ([]string, error) {
	s3 := env.S3(ctx)
	screenPath := path.Join(constants.S3Prefix, docID, "screenshots")
	screenKeys, err := s3.ListFirstTwo(s3.Bucket, screenPath)
	if err != nil {
		return nil, fmt.Errorf("error getting last screenshot: %s", err)
	}

	darkUrl, err := s3.GetPresignedUrl(
		s3.Bucket,
		screenKeys[0],
		screenshotUrlExpireDuration,
	)
	if err != nil {
		return nil, err
	}

	lightUrl, err := s3.GetPresignedUrl(
		s3.Bucket,
		screenKeys[1],
		screenshotUrlExpireDuration,
	)
	if err != nil {
		return nil, err
	}

	return []string{darkUrl, lightUrl}, nil
}
