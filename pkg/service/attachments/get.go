package attachments

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func GetForUser(ctx context.Context, userID string, ids []string) ([]*models.DocumentAttachment, error) {
	if len(ids) == 0 {
		return []*models.DocumentAttachment{}, nil
	}

	docAttTbl := env.Query(ctx).DocumentAttachment

	attachments, err := docAttTbl.Where(
		docAttTbl.UserID.Eq(userID),
		docAttTbl.ID.In(ids...),
	).Find()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		found := false
		for _, attchmnt := range attachments {
			if attchmnt.ID == id {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("could not find attachment with id %s", id)
		}
	}

	return attachments, nil
}
