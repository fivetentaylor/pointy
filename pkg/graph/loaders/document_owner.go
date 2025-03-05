package loaders

import (
	"context"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"gorm.io/gen/field"
)

type documentOwnerReader struct{}

// getDocumentOwners implements a batch function that can retrieve many document owners by document ID,
// for use in a dataloader
func (d *documentOwnerReader) getDocumentOwners(ctx context.Context, documentIDs []string) ([]*models.User, []error) {
	log := env.Log(ctx)
	userTbl := env.Query(ctx).User
	docAccessTbl := env.Query(ctx).DocumentAccess

	type UserWithDocID struct {
		models.User
		DocumentID string
	}

	var usersWithDocIDs []UserWithDocID

	fields := []field.Expr{
		userTbl.ID,
		userTbl.Name,
		userTbl.Email,
		userTbl.CreatedAt,
		userTbl.UpdatedAt,
		userTbl.Picture,
		userTbl.Admin,
		userTbl.DisplayName,
		docAccessTbl.DocumentID,
	}

	err := userTbl.WithContext(ctx).
		Select(fields...).
		Where(docAccessTbl.DocumentID.In(documentIDs...)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		LeftJoin(docAccessTbl, userTbl.ID.EqCol(docAccessTbl.UserID)).
		Scan(&usersWithDocIDs)

	if err != nil {
		log.Errorf("error getting document owners: %s", err)
		return nil, []error{err}
	}

	userMap := make(map[string]*models.User)
	for _, userWithDocID := range usersWithDocIDs {
		userMap[userWithDocID.DocumentID] = &userWithDocID.User
	}

	// Ensure the result slice is in the same order as the input keys
	out := make([]*models.User, len(documentIDs))
	for i, id := range documentIDs {
		out[i] = userMap[id]
	}

	return out, nil
}

// GetDocumentOwner returns a single document owner by document ID efficiently
func GetDocumentOwner(ctx context.Context, documentID string) (*models.User, error) {
	loaders := For(ctx)
	return loaders.DocumentOwnerLoader.Load(ctx, documentID)
}

// GetDocumentOwners returns many document owners by document IDs efficiently
func GetDocumentOwners(ctx context.Context, documentIDs []string) ([]*models.User, error) {
	loaders := For(ctx)
	return loaders.DocumentOwnerLoader.LoadAll(ctx, documentIDs)
}
