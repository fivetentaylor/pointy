package loaders

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/env"
)

type DocumentAccessInput struct {
	UserID     string
	DocumentID string
}

// getDocumentAccesss implements a batch function that can retrieve many document owners by document ID,
// for use in a dataloader
func getDocumentAccesss(ctx context.Context, inputs []DocumentAccessInput) ([]string, []error) {
	docAccessTbl := env.Query(ctx).DocumentAccess

	userIDs := make([]string, len(inputs))
	documentIDs := make([]string, len(inputs))
	for i, docID := range inputs {
		userIDs[i] = docID.UserID
		documentIDs[i] = docID.DocumentID
	}

	accesses, err := docAccessTbl.
		Select(docAccessTbl.AccessLevel, docAccessTbl.UserID, docAccessTbl.DocumentID).
		Where(docAccessTbl.UserID.In(userIDs...)).
		Where(docAccessTbl.DocumentID.In(documentIDs...)).
		Find()
	if err != nil {
		return nil, []error{err}
	}

	out := make([]string, len(inputs))
	for i, docID := range inputs {
		for _, access := range accesses {
			if access.DocumentID == docID.DocumentID && access.UserID == docID.UserID {
				out[i] = access.AccessLevel
				break
			}
		}
	}

	return out, nil
}

// GetDocumentAccess returns a single document owner by document ID efficiently
func GetDocumentAccess(ctx context.Context, documentID, userID string) (string, error) {
	loaders := For(ctx)
	return loaders.DocumentAccessLoader.Load(ctx, DocumentAccessInput{UserID: userID, DocumentID: documentID})
}
