package sharing

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/timeline"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func UpdateLink(ctx context.Context, userId string, link *models.SharedDocumentLink, isActive bool) (*models.SharedDocumentLink, error) {
	shareLinkTbl := env.Query(ctx).SharedDocumentLink

	_, err := shareLinkTbl.
		Where(shareLinkTbl.ID.Eq(link.ID)).
		Update(shareLinkTbl.IsActive, isActive)
	if err != nil {
		log.Errorf("error updating share link: %s", err)
		return nil, fmt.Errorf("sorry, we could not update your share link")
	}

	if !isActive && link.IsActive {
		event := &dynamo.TimelineEvent{
			DocID:  link.DocumentID,
			UserID: userId,
			Event: &models.TimelineEventPayload{
				Payload: &models.TimelineEventPayload_AccessChange{
					AccessChange: &models.TimelineAccessChangeV1{
						Action: models.TimelineAccessChangeAction_REMOVE_ACTION,
						UserIdentifiers: []string{
							link.InviteeEmail,
						},
					},
				},
			},
		}

		err = timeline.CreateTimelineEvent(ctx, event)
		if err != nil {
			log.Errorf("error creating timeline event: %s", err)
		}
	}

	link.IsActive = isActive

	return link, nil
}
