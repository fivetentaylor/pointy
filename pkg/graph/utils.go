package graph

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func checkIfOwner(ctx context.Context, docID string) error {
	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Errorf("error getting current user: %s", err)
		return fmt.Errorf("please login")
	}

	docAccessTbl := env.Query(ctx).DocumentAccess

	// Check if the current user is the owner of the document
	count, err := docAccessTbl.
		Where(docAccessTbl.DocumentID.Eq(docID)).
		Where(docAccessTbl.UserID.Eq(currentUser.Id)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		Count()
	if err != nil {
		log.Errorf("error checking document access: %s", err)
		return fmt.Errorf("internal error")
	}

	if count == 0 {
		log.Errorf("current user is not the owner of the document")
		return fmt.Errorf("unauthorized")
	}

	return nil
}
