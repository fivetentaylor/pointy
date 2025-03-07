package pubsub

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func loadDocument(ctx context.Context, docID string) (*models.Document, error) {
	documentTbl := env.Query(ctx).Document

	return documentTbl.
		Where(documentTbl.ID.Eq(docID)).
		First()
}

func ListenForNewDocuments(ctx context.Context, ch chan *models.Document, userID string) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)

	pubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.ChannelNewDocsFormat, userID))
	incoming := pubsub.Channel()

	defer func() {
		log.Info("closing new document listener", "userID", userID)
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		close(ch)
	}()

	for {
		select {
		case msg := <-incoming:
			newDocID := msg.Payload
			doc, err := loadDocument(ctx, newDocID)
			if err != nil {
				return
			}

			ch <- doc
		case <-ctx.Done():
			return
		}
	}
}

func ListenForDocumentUpdates(ctx context.Context, ch chan *models.Document, docID, userID string) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)

	pubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.DocUpsertChanFormat, docID))
	incoming := pubsub.Channel()

	userPubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.UserDocUpdatesChanFormat, docID, userID))
	userIncoming := userPubsub.Channel()

	defer func() {
		log.Info("closing document listener", "docID", docID)
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		userPubsub.Unsubscribe(ctx)
		userPubsub.Close()
		close(ch)
	}()

	for {
		log.Debug("waiting for document updates")

		doc, err := loadDocument(ctx, docID)
		if err != nil {
			return
		}

		ch <- doc

		select {
		case <-userIncoming:
			doc, err := loadDocument(ctx, docID)
			if err != nil {
				return
			}

			ch <- doc
		case <-incoming:
			doc, err := loadDocument(ctx, docID)
			if err != nil {
				return
			}

			ch <- doc
		case <-ctx.Done():
			return
		}
	}
}

func PublishDocument(ctx context.Context, docID string) error {
	log := env.Log(ctx)
	log.Info("publishing document", "docID", docID)
	key := fmt.Sprintf(constants.DocUpsertChanFormat, docID)
	return env.Redis(ctx).Publish(ctx, key, docID).Err()
}

func PublishUserDocument(ctx context.Context, docID, userID string) error {
	key := fmt.Sprintf(constants.UserDocUpdatesChanFormat, docID, userID)
	return env.Redis(ctx).Publish(ctx, key, docID).Err()
}

func PublishNewDocument(ctx context.Context, userID, docID string) error {
	key := fmt.Sprintf(constants.ChannelNewDocsFormat, userID)
	return env.Redis(ctx).Publish(ctx, key, docID).Err()
}
