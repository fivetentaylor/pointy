package sharing

import (
	"context"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/timeline"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func AddTimelineJoin(ctx context.Context, doc *models.Document, userID string) error {
	log := env.Log(ctx)

	event := &dynamo.TimelineEvent{
		DocID:  doc.ID,
		UserID: userID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Join{
				Join: &models.TimelineJoinV1{
					Action: "join",
				},
			},
		},
	}

	err := timeline.CreateTimelineEvent(ctx, event)
	if err != nil {
		log.Errorf("error creating timeline event: %s", err)
		return err
	}

	return nil
}
