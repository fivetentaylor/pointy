package sharing

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/pubsub"
	"github.com/fivetentaylor/pointy/pkg/service/timeline"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func UnshareDoc(
	ctx context.Context,
	userID string,
	document *models.Document,
	editorID string,
) error {
	log := env.Log(ctx)
	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		docAccessTbl := tx.DocumentAccess

		accessRow, err := docAccessTbl.
			Where(docAccessTbl.DocumentID.Eq(document.ID)).
			Where(docAccessTbl.UserID.Eq(editorID)).
			First()

		if err != nil {
			log.Errorf("error getting access row: %s", err)
			return err
		}

		_, err = docAccessTbl.Delete(accessRow)

		return err
	})
	if err != nil {
		return fmt.Errorf("failed to unshare document: %s", err)
	}

	editor, err := env.Query(ctx).User.Where(env.Query(ctx).User.ID.Eq(editorID)).First()
	if err != nil {
		return fmt.Errorf("failed to unshare document: %s", err)
	}

	err = pubsub.PublishDocument(ctx, document.ID)
	if err != nil {
		log.Errorf("error publishing document: %s", err)
	}

	event := &dynamo.TimelineEvent{
		DocID:  document.ID,
		UserID: userID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_AccessChange{
				AccessChange: &models.TimelineAccessChangeV1{
					Action: models.TimelineAccessChangeAction_REMOVE_ACTION,
					UserIdentifiers: []string{
						editor.Name,
					},
				},
			},
		},
	}

	err = timeline.CreateTimelineEvent(ctx, event)
	if err != nil {
		log.Errorf("error creating timeline event: %s", err)
	}

	return nil
}
