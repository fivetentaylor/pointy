package dag

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

const documentStoreKey = "_docStore"                    // docId
const documentKeyFormat = "_document:%s:%s"             // docId, authorId
const documentAtMessageKeyFormat = "_document:%s:%s:%s" // docId, authorId, messageId
const realtimeKeyFormat = "_document_rt:%s"             // docId

func GetDocument(ctx context.Context, docId, authorId string) (*v3.Rogue, error) {
	document, err := GetStateKey[*v3.Rogue](ctx, fmt.Sprintf(documentKeyFormat, docId, authorId))
	if err != nil {
		return nil, err
	}
	if document == nil {
		docStore, err := GetDocStore(ctx)
		if err != nil {
			return nil, err
		}

		_, document, err = docStore.GetCurrentDoc(ctx, docId)
		if err != nil {
			return nil, err
		}

		document.Author = authorId

		SetStateKey(ctx, fmt.Sprintf(documentKeyFormat, docId, authorId), document)
	}

	if document == nil {
		return nil, fmt.Errorf("no document found")
	}

	return document, err
}

func GetDocumentAtMessageID(ctx context.Context, docId, authorId, threadId, messageId string) (*v3.Rogue, error) {
	document, err := GetStateKey[*v3.Rogue](ctx, fmt.Sprintf(documentAtMessageKeyFormat, docId, authorId, messageId))
	if err != nil {
		return nil, err
	}
	if document == nil {
		msg, err := env.Dynamo(ctx).GetAiThreadMessage(threadId, messageId)
		if err != nil {
			return nil, fmt.Errorf("could not get ai thread %s message %s: %w", threadId, messageId, err)
		}

		document, err = GetDocumentAtMessage(ctx, docId, authorId, msg)
		if err != nil {
			return nil, fmt.Errorf("could not get document at message: %w", err)
		}

		document.Author = authorId
	}

	if document == nil {
		return nil, fmt.Errorf("no document found")
	}

	return document, err
}

func GetDocumentAtMessage(ctx context.Context, docId, authorId string, message *dynamo.Message) (*v3.Rogue, error) {
	document, err := GetStateKey[*v3.Rogue](ctx, fmt.Sprintf(documentAtMessageKeyFormat, docId, authorId, message.MessageID))
	if err != nil {
		return nil, err
	}
	if document == nil {
		docStore, err := GetDocStore(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get document store: %w", err)
		}

		_, document, err = docStore.GetCurrentDoc(ctx, docId)
		if err != nil {
			return nil, fmt.Errorf("could not get current doc %s document: %w", docId, err)
		}

		ca, err := v3.ParseContentAddress(message.MessageMetadata.ContentAddress)
		if err != nil {
			return nil, fmt.Errorf("could not parse content address %q: %w", message.MessageMetadata.ContentAddress, err)
		}

		document, err = document.GetOldRogue(ca)
		if err != nil {
			return nil, fmt.Errorf("could not get old rogue: %w", err)
		}

		document.Author = authorId
	}

	if document == nil {
		return nil, fmt.Errorf("no document found")
	}

	return document, err
}

func GetDocStore(ctx context.Context) (*rogue.DocStore, error) {
	ds, err := GetStateKey[*rogue.DocStore](ctx, documentStoreKey)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		ds = rogue.NewDocStore(env.S3(ctx), env.Query(ctx), env.Redis(ctx))
		SetStateKey(ctx, documentStoreKey, ds)
	}

	return ds, nil
}

func GetDocumentRealtimeConnection(ctx context.Context, docId string) (*rogue.Realtime, error) {
	rt, err := GetStateKey[*rogue.Realtime](ctx, fmt.Sprintf(realtimeKeyFormat, docId))
	if err != nil {
		return nil, err
	}
	if rt == nil {
		rt = rogue.NewRealtime(
			env.Redis(ctx),
			env.Query(ctx),
			docId,
			constants.RevisoUserID,
			"Reviso",
			constants.ColorIndigo,
		)

		SetStateKey(ctx, fmt.Sprintf(realtimeKeyFormat, docId), rt)
	}

	if rt == nil {
		return nil, fmt.Errorf("no document realtime connection found")
	}

	return rt, err
}
