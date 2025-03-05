package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
)

const avatarSize = 192
const avatarUrlDuration = 60 * time.Minute

func AvatarUrlForUser(ctx context.Context, user *models.User) (*string, error) {
	if user.Picture == nil {
		return nil, nil
	}

	pictureURL := *user.Picture

	parsedURL, err := url.Parse(pictureURL)
	if err != nil {
		log.Errorf("error parsing picture url %s: %s", pictureURL, err)
		return nil, err
	}

	if parsedURL.Scheme == "" {
		s3 := env.S3(ctx)
		presignedUrl, err := s3.GetPresignedUrl(s3.Bucket, pictureURL, avatarUrlDuration)
		if err != nil {
			log.Errorf("error getting presigned url: %s", err)
			return nil, err
		}
		return &presignedUrl, nil
	}

	return &pictureURL, nil
}

func UpdateUserAvatar(ctx context.Context, user *models.User, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	var resizedImage image.Image

	if width != height {
		resizedImage = ResizeImage(CropToCenterSquare(img), avatarSize, avatarSize)
	} else {
		resizedImage = ResizeImage(img, avatarSize, avatarSize)
	}

	var outputBuffer bytes.Buffer
	err := png.Encode(&outputBuffer, resizedImage)
	if err != nil {
		return fmt.Errorf("error encoding image: %w", err)
	}

	s3 := env.S3(ctx)
	s3key := user.AvatarS3Key()
	err = s3.PutObject(s3.Bucket, s3key, "image/png", outputBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("error uploading image: %w", err)
	}

	userTbl := env.Query(ctx).User
	user.Picture = &s3key
	err = userTbl.Save(user)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}
