package attachments

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func ListForDocumentAndUser(ctx context.Context, docID, userID string) ([]*models.DocumentAttachment, error) {
	docAttTbl := env.Query(ctx).DocumentAttachment

	attachments, err := docAttTbl.Where(
		docAttTbl.DocumentID.Eq(docID),
		docAttTbl.UserID.Eq(userID),
	).Find()
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func ListForUser(ctx context.Context, userID string) ([]*models.DocumentAttachment, error) {
	docAttTbl := env.Query(ctx).DocumentAttachment

	attachments, err := docAttTbl.Where(
		docAttTbl.UserID.Eq(userID),
	).Find()
	if err != nil {
		return nil, err
	}

	return attachments, nil
}
