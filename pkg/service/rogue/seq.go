package rogue

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
)

type ErrorPrefixNotFound struct {
	Prefix string
}

func (e ErrorPrefixNotFound) Error() string {
	return fmt.Sprintf("no prefix: %s", e.Prefix)
}

func GetLastS3Seq(ctx context.Context, docID string) (int64, error) {
	snapPath := DocS3Path(docID)
	snapKey, err := ListFirst(ctx, snapPath)
	if err != nil {
		if errors.As(err, &ErrorPrefixNotFound{}) {
			return 0, nil
		} else {
			return 0, err
		}
	}

	seq, err := ExtractSeqFromS3Path(snapKey)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

func ExtractSeqFromS3Path(s3Path string) (int64, error) {
	parts := strings.Split(s3Path, "/")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid S3 path: %s", s3Path)
	}
	invSeq := parts[3]

	return revertInvertedSeq(invSeq)
}

func ListFirst(ctx context.Context, prefix string) (string, error) {
	s3Client := env.S3(ctx)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3Client.Bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(1),
	}

	result, err := s3Client.Client.ListObjectsV2(input)
	if err != nil {
		log.Errorf("failed to list objects: %s", err)
		return "", err
	}

	if len(result.Contents) == 0 {
		return "", ErrorPrefixNotFound{Prefix: prefix}
	}

	return *result.Contents[0].Key, nil
}
