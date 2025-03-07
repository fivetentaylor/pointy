package document

import (
	"context"
	"fmt"
	"time"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/pubsub"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func Create(ctx context.Context, userID string) (*models.Document, error) {
	log := env.SLog(ctx)
	doc := &models.Document{Title: constants.DefaultDocumentTitle}

	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		documentTbl := tx.Document
		docAccessTbl := tx.DocumentAccess

		err := documentTbl.Omit(documentTbl.RootParentID).Create(doc)
		if err != nil {
			log.Error("error creating db document", "error", err)
			return err
		}

		err = docAccessTbl.Create(&models.DocumentAccess{
			UserID:         userID,
			DocumentID:     doc.ID,
			AccessLevel:    constants.DefaultDocumentRole,
			LastAccessedAt: time.Now(),
		})
		if err != nil {
			log.Error("error creating document access", "error", err)
			return err
		}

		return AddTimelineCreate(ctx, doc, userID)
	})
	if err != nil {
		log.Error("error creating db document", "error", err)
		return nil, fmt.Errorf("sorry, we could not create your document")
	}

	err = pubsub.PublishNewDocument(ctx, userID, doc.ID)
	if err != nil {
		log.Error("error publishing document", "error", err)
	}

	log.Info("document created", "doc_id", doc.ID, "event", "document_created")
	return doc, nil
}

func AddTimelineCreate(ctx context.Context, doc *models.Document, userID string) error {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)

	event := &dynamo.TimelineEvent{
		DocID:  doc.ID,
		UserID: userID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Join{
				Join: &models.TimelineJoinV1{
					Action: "create",
				},
			},
		},
	}

	err := dydb.CreateTimelineEvent(event)
	if err != nil {
		log.Errorf("error creating timeline event: %s", err)
		return err
	}

	return nil
}

func CreateCustom(ctx context.Context, userID string, doc *models.Document) (*models.Document, error) {
	log := env.SLog(ctx)
	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		documentTbl := tx.Document
		docAccessTbl := tx.DocumentAccess

		if doc.RootParentID == "" {
			err := documentTbl.Omit(documentTbl.RootParentID).Create(doc)
			if err != nil {
				log.Error("error creating db document", "error", err)
				return err
			}
		} else {
			err := documentTbl.Create(doc)
			if err != nil {
				log.Error("error creating db document", "error", err)
				return err
			}
		}

		var err error
		doc, err = documentTbl.Where(documentTbl.ID.Eq(doc.ID)).First()
		if err != nil {
			log.Error("error getting document", "error", err)
			return err
		}

		err = docAccessTbl.Create(&models.DocumentAccess{
			UserID:         userID,
			DocumentID:     doc.ID,
			AccessLevel:    constants.DefaultDocumentRole,
			LastAccessedAt: time.Now(),
		})
		if err != nil {
			log.Error("error creating document access", "error", err)
			return err
		}

		return AddTimelineCreate(ctx, doc, userID)
	})
	if err != nil {
		log.Error("error creating db document", "error", err)
		return nil, fmt.Errorf("sorry, we could not create your document")
	}

	err = pubsub.PublishNewDocument(ctx, userID, doc.ID)
	if err != nil {
		log.Error("error publishing document", "error", err)
	}

	log.Info("document created", "doc_id", doc.ID, "event", "document_created")
	return doc, nil
}
