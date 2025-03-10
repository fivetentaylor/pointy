package rogue

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/env"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func HTMLSelection(ctx context.Context, docId, start, end string) (string, error) {
	log := env.Log(ctx)
	doc, err := CurrentDocument(ctx, docId)
	if err != nil {
		log.Errorf("error getting current document: %s", err)
		return "", err
	}

	startId, err := v3.ParseID(start)
	if err != nil {
		log.Errorf("error parsing start id: %s", err)
		return "", err
	}

	endId, err := v3.ParseID(end)
	if err != nil {
		log.Errorf("error parsing end id: %s", err)
		return "", err
	}

	return doc.GetHtml(startId, endId, false, false)
}
