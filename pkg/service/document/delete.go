package document

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func Delete(ctx context.Context, documentId string) error {
	id := documentId
	docAccessTbl := env.Query(ctx).DocumentAccess

	users, err := docAccessTbl.Where(docAccessTbl.DocumentID.Eq(id)).Select(docAccessTbl.UserID).Find()
	if err != nil {
		return fmt.Errorf("could not find document to delete: %s", err)
	}

	for _, user := range users {
		_, err := DeleteMessaging(ctx, documentId, user.UserID)
		if err != nil {
			return err
		}
	}

	err = DeleteStorage(ctx, documentId)
	if err != nil {
		return err
	}

	documentTbl := env.Query(ctx).Document

	_, err = docAccessTbl.
		Where(docAccessTbl.DocumentID.Eq(id)).
		Delete()
	if err != nil {
		return fmt.Errorf("could not find document to delete: %s", err)
	}

	_, err = documentTbl.
		Where(documentTbl.ID.Eq(id)).
		Delete()
	if err != nil {
		return fmt.Errorf("could not delete document: %s", err)
	}

	return nil
}

func DeleteStorage(ctx context.Context, documentId string) error {
	s3 := env.S3(ctx)

	err := s3.DeletePrefix(s3.Bucket, fmt.Sprintf(constants.DocumentSnapshotPrefix, constants.S3Prefix, documentId))
	if err != nil {
		return err
	}

	err = s3.DeletePrefix(s3.Bucket, fmt.Sprintf(constants.ConvosPrefix, documentId))
	if err != nil {
		return err
	}

	err = s3.DeletePrefix(s3.Bucket, fmt.Sprintf(constants.LogsPrefix, documentId))
	if err != nil {
		return err
	}

	return nil
}

func DeleteMessaging(ctx context.Context, documentId, userID string) ([]*dynamo.Thread, error) {
	dydb := env.Dynamo(ctx)

	threads, err := dydb.DeleteThreadsForDoc(documentId, userID)
	if err != nil {
		return nil, err
	}

	err = dydb.DeleteDocumentChannels(documentId)
	if err != nil {
		return threads, err
	}

	return threads, nil
}
