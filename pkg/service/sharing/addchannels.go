package sharing

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func AddChannels(ctx context.Context, doc *models.Document, userID string) error {
	docAccessTbl := env.Query(ctx).DocumentAccess

	users, err := docAccessTbl.Where(docAccessTbl.DocumentID.Eq(doc.ID)).Select(docAccessTbl.UserID).Find()
	if err != nil {
		log.Errorf("error getting user ids: %s", err)
		return fmt.Errorf("sorry, we could not share your document")
	}

	userIds := make([]string, len(users))
	for i, user := range users {
		userIds[i] = user.UserID
		log.Infof("creating channel for: %s", user.UserID)
	}

	return nil
}
