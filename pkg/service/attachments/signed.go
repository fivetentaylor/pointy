package attachments

import (
	"context"
	"fmt"
	"time"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/graph/model"
)

func GetSignedURL(ctx context.Context, userID, attachmentID, filename string) (*model.SignedImageURL, error) {
	s3 := env.S3(ctx)

	docAttTbl := env.Query(ctx).DocumentAttachment

	attachment, err := docAttTbl.Where(
		docAttTbl.UserID.Eq(userID),
		docAttTbl.ID.Eq(attachmentID),
	).First()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf(constants.DocumentAttachmentFileKey, attachment.S3ID, filename)

	url, err := s3.GetPresignedUrl(s3.ImagesBucket, key, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	return &model.SignedImageURL{
		URL:       url,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}, nil
}
