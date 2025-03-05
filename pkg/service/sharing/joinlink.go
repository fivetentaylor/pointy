package sharing

import (
	"context"
	"fmt"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
)

func JoinDoc(ctx context.Context, docId string, userID string, accessLevel string) (*models.Document, error) {
	log := env.Log(ctx)
	documentTbl := env.Query(ctx).Document

	doc, err := documentTbl.
		Where(documentTbl.ID.Eq(docId)).
		First()
	if err != nil {
		return nil, fmt.Errorf("could not find document: %s", err)
	}

	existingAccess, err := env.Query(ctx).DocumentAccess.
		Where(
			env.Query(ctx).DocumentAccess.DocumentID.Eq(doc.ID),
			env.Query(ctx).DocumentAccess.UserID.Eq(userID),
		).First()
	if err == nil && existingAccess != nil {
		log.Warn("user already joined document", "userID", userID, "docID", doc.ID)
		return doc, nil
	}

	log.Info("joining document", "userID", userID, "docID", doc.ID)
	err = env.Query(ctx).Transaction(func(tx *query.Query) error {
		docAccessTbl := tx.DocumentAccess

		err = docAccessTbl.Create(&models.DocumentAccess{
			UserID:      userID,
			DocumentID:  doc.ID,
			AccessLevel: accessLevel,
		})
		return err
	})
	if err != nil {
		log.Errorf("error joining document: %s", err)
		return nil, fmt.Errorf("failed to join document: %s", err)
	}

	return doc, nil
}
