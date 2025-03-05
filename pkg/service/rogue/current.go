package rogue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func CurrentDocument(ctx context.Context, docId string) (*v3.Rogue, error) {
	_, doc, err := CurrentDocumentAndSequence(ctx, docId)
	return doc, err
}

func CurrentDocumentAndSequence(ctx context.Context, docId string) (int64, *v3.Rogue, error) {
	log := env.Log(ctx)
	seq, doc, err := LastSnapshot(ctx, docId)
	if err != nil {
		log.Errorf("[doc %s] error getting last snapshot: %s", docId, err)
		return 0, nil, err
	}

	ops, err := getDeltaLog(ctx, docId, 0, -1)
	if err != nil {
		return 0, nil, err
	}

	if len(ops) == 0 {
		return seq, doc, nil
	}

	var maxSeq int64
	for _, z := range ops {
		op := z.Member.(string)
		if strings.HasPrefix(op, "deleted") {
			continue
		}

		var msg v3.Message
		err = json.Unmarshal([]byte(op), &msg)
		if err != nil {
			log.Errorf("[doc %s] error unmarshalling op: %s", docId, err)
			return 0, nil, err
		}

		_, err := doc.MergeOp(msg.Op)
		if err != nil {
			log.Errorf("[doc %s] (ignored) Error merging op %#v: %s", docId, msg.Op, err)
		}

		maxSeq = int64(z.Score)
	}

	return maxSeq, doc, nil

}

func LastSnapshot(ctx context.Context, docId string) (int64, *v3.Rogue, error) {
	s3Client := env.S3(ctx)

	prefix := fmt.Sprintf(constants.DocumentSnapshotPrefix, constants.S3Prefix, docId)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3Client.Bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(1),
	}

	result, err := s3Client.Client.ListObjectsV2(input)
	if err != nil {
		log.Errorf("failed to list objects: %s", err)
		return 0, nil, err
	}

	if len(result.Contents) == 0 {
		// No snapshots, return empty document
		return 0, v3.NewRogueForQuill("s"), nil
	}

	snapKey := *result.Contents[0].Key
	snap, err := s3Client.GetObject(s3Client.Bucket, snapKey)
	if err != nil {
		return 0, nil, err
	}

	var doc v3.Rogue
	err = json.Unmarshal(snap, &doc)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling snapshot: %w", err)
	}

	seq, err := ExtractSeqFromS3Path(snapKey)
	if err != nil {
		return 0, nil, err
	}

	return seq, &doc, nil
}

func getDeltaLog(ctx context.Context, docId string, lowScore, highScore int) ([]redis.Z, error) {
	rds := env.Redis(ctx)
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docId)

	max := fmt.Sprintf("%d", highScore)
	if highScore == -1 {
		max = "+inf"
	}

	// ZRangeByScore gets all elements in the sorted set within the specified score range
	events, err := rds.ZRangeByScoreWithScores(ctx, eventsZSetKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", lowScore),
		Max: max,
	}).Result()

	if err != nil {
		return nil, err
	}

	log.Infof("GetDeltaLog(%q, %d, %d) -> %d events", docId, lowScore, highScore, len(events))

	return events, nil
}

func revertInvertedSeq(invSeq string) (int64, error) {
	invertedSeq, err := strconv.ParseInt(invSeq, 10, 64)
	if err != nil {
		return 0, err
	}
	return constants.MaxSeqValue - invertedSeq, nil
}
