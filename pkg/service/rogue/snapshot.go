package rogue

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func SnapshotDoc(ctx context.Context, docId string) error {
	seq, doc, err := CurrentDocumentAndSequence(ctx, docId)
	if err != nil {
		return err
	}

	err = SaveDocToS3(ctx, docId, seq, doc)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	err = cleanupDeltaLog(ctx, docId, seq)
	if err != nil {
		return err
	}

	return nil
}

func cleanupDeltaLog(ctx context.Context, docID string, highScore int64) error {
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	rds := env.Redis(ctx)
	_, err := rds.ZRemRangeByScore(ctx, eventsZSetKey, "0", fmt.Sprintf("%d", highScore)).Result()
	if err != nil {
		return err
	}

	return nil
}
