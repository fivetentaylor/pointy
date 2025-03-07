package attachments

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func Upload(ctx context.Context, file graphql.Upload, docID, userID string) (*models.DocumentAttachment, error) {
	s3 := env.S3(ctx)
	log := env.SLog(ctx)

	id := uuid.NewString()
	key := fmt.Sprintf(constants.DocumentAttachmentOriginalKey, id)

	data, err := io.ReadAll(file.File)
	if err != nil {
		log.Error("error reading file", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
		return nil, fmt.Errorf("internal error")
	}

	url := strings.TrimSpace(string(data))
	data = []byte(url)

	err = s3.PutObject(s3.Bucket, key, file.ContentType, data)
	if err != nil {
		log.Error("error uploading image", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
		return nil, fmt.Errorf("internal error")
	}

	attachment := &models.DocumentAttachment{
		UserID:      userID,
		DocumentID:  docID,
		Filename:    file.Filename,
		Size:        file.Size,
		S3ID:        id,
		ContentType: file.ContentType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Special case for URL
	if file.ContentType == "text/url" {
		attachment.Filename = url
		// Read the URL to make sure we can access it
		_, err := ExtractedText(ctx, attachment)
		if err != nil {
			log.Error("error extracting text", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
			return nil, fmt.Errorf("Could not read URL %s. Some domain block services like Reviso from accessing their content.", url)
		}
	}

	log.Info("creating db record for attachment", slog.Any("attachment", attachment))

	docAttTb := env.Query(ctx).DocumentAttachment
	err = docAttTb.Create(attachment)
	if err != nil {
		log.Error("error creating db record for attachment", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
		return nil, fmt.Errorf("internal error")
	}

	return attachment, nil
}
