package document

import (
	"context"
	"fmt"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/pubsub"
	"github.com/teamreviso/code/pkg/service/timeline"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func UpdateDocument(ctx context.Context, doc *models.Document, userID string, title *string, isPublic *bool) (*models.Document, error) {
	log := env.Log(ctx)

	previousTitle := doc.Title
	if title != nil {
		doc.Title = *title
	}
	previousIsPublic := doc.IsPublic
	if isPublic != nil {
		doc.IsPublic = *isPublic
	}

	// we have to do this because of how gorm works
	// just using .Save falls back to an Upsert that doens't like a default TRUE on a column
	// see GORM note: "NOTE When updating with struct, GORM will only update non-zero fields. You might want to use map to update attributes or use Select to specify fields to update" via https://gorm.io/docs/update.html
	err := env.RawDB(ctx).Model(&doc).Updates(map[string]interface{}{"title": doc.Title, "is_public": doc.IsPublic}).Error

	if err != nil {
		log.Errorf("error updating document: %s", err)
		return nil, fmt.Errorf("sorry, we could not update your document")
	}

	err = pubsub.PublishDocument(ctx, doc.ID)
	if err != nil {
		log.Errorf("error publishing document: %s", err)
	}

	if doc.Title != previousTitle {
		event := &dynamo.TimelineEvent{
			DocID:  doc.ID,
			UserID: userID,
			Event: &models.TimelineEventPayload{
				Payload: &models.TimelineEventPayload_AttributeChange{
					AttributeChange: &models.TimelineAttributeChangedV1{
						Attribute: "title",
						OldValue:  previousTitle,
						NewValue:  doc.Title,
					},
				},
			},
		}

		err = timeline.CreateTimelineEvent(ctx, event)
		if err != nil {
			log.Errorf("error creating timeline event: %s", err)
		}
	}

	if doc.IsPublic != previousIsPublic {
		event := &dynamo.TimelineEvent{
			DocID:  doc.ID,
			UserID: userID,
			Event: &models.TimelineEventPayload{
				Payload: &models.TimelineEventPayload_AttributeChange{
					AttributeChange: &models.TimelineAttributeChangedV1{
						Attribute: "is_public",
						OldValue:  fmt.Sprintf("%t", previousIsPublic),
						NewValue:  fmt.Sprintf("%t", doc.IsPublic),
					},
				},
			},
		}

		err = timeline.CreateTimelineEvent(ctx, event)
		if err != nil {
			log.Errorf("error creating timeline event: %s", err)
		}
	}

	return doc, nil
}
