package rogue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/pkg/storage/s3"
	rogueV3 "github.com/teamreviso/code/rogue/v3"
)

const checkpointInterval = 1000

type DocStore struct {
	S3    *s3.S3
	Query *query.Query
	Redis *redis.Client
}

func NewDocStore(s3 *s3.S3, query *query.Query, redis *redis.Client) *DocStore {
	return &DocStore{
		S3:    s3,
		Query: query,
		Redis: redis,
	}
}

func (ds *DocStore) ConvertS3IfPresent(ctx context.Context, docID string) error {
	eventsPrefix := path.Join("v2", docID)
	deltas, err := ds.S3.GetAllObjects(ds.S3.Bucket, eventsPrefix)
	if err != nil {
		return err
	}

	for _, delta := range deltas {
		var msg rogueV3.Message
		err = json.Unmarshal(delta, &msg)
		if err != nil {
			return err
		}

		_, err := ds.AddUpdate(ctx, docID, msg.Op)
		if err != nil {
			return err
		}
	}

	err = ds.S3.DeleteAll(ds.S3.Bucket, eventsPrefix)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DocStore) LoadLastCheckpoint(docID string) ([]byte, error) {
	return nil, nil
}

func (ds *DocStore) AddUpdate(ctx context.Context, docID string, op rogueV3.Op) (int64, error) {
	bytes, err := json.Marshal(op)
	if err != nil {
		return 0, err
	}

	seq, err := ds.AddDeltaLog(context.Background(), docID, string(bytes))
	if err != nil {
		return 0, err
	}

	return seq, err
}

func (ds *DocStore) DocOperationsKey(docID string) string {
	return fmt.Sprintf("%s/ops", docID)
}

func (ds *DocStore) OperationKey(docID string, op rogueV3.Op) string {
	return fmt.Sprintf("%s/%020d_%s", ds.DocOperationsKey(docID), op.GetID().Seq, op.GetID().Author)
}

func (ds *DocStore) Checkpoint(docID string) error {
	return nil
}

const luaScript = `if redis.call("EXISTS", KEYS[1]) == 0 then
	redis.call("SET", KEYS[1], ARGV[1])
end
return redis.call("INCR", KEYS[1])`

func (ds *DocStore) AddDeltaLog(ctx context.Context, docID string, eventData string) (int64, error) {
	var err error
	rdb := ds.Redis

	// Define keys for the counter and the sorted set of event data specific to the docID
	counterKey := fmt.Sprintf(constants.DocCounterKeyFormat, docID)
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	// Get the next sequence number for the docID
	var score int64
	// Check if the counter key exists
	getCmd := rdb.Get(ctx, counterKey)
	if getCmd.Err() == redis.Nil {
		// if the counter key does not exist, create it from the s3 sequence number
		s3seq, err := ds.GetLastS3Seq(docID)
		if err != nil {
			return 0, err
		}

		// Attempt to set the counter key, if it does not exist (i.e it has been created since we called GET)
		seqIDCmd := rdb.Eval(ctx, luaScript, []string{counterKey}, s3seq)
		if seqIDCmd.Err() != nil {
			return 0, seqIDCmd.Err()
		}

		score = seqIDCmd.Val().(int64)
	} else if getCmd.Err() != nil {
		return 0, getCmd.Err()
	} else {
		// if the counter key exists, increment it
		seqIDCmd := rdb.Incr(ctx, counterKey)
		err = seqIDCmd.Err()
		if err != nil {
			return 0, err
		}
		score = seqIDCmd.Val()
	}

	// Add the event data to the sorted set
	setCmd := rdb.ZAdd(ctx, eventsZSetKey, redis.Z{
		Score:  float64(score),
		Member: eventData,
	})
	if setCmd.Err() != nil {
		return 0, setCmd.Err()
	}

	return score, nil
}

func (ds *DocStore) CleanupDeltaLog(ctx context.Context, docID string, highScore int64) error {
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	_, err := ds.Redis.ZRemRangeByScore(ctx, eventsZSetKey, "0", fmt.Sprintf("%d", highScore)).Result()
	if err != nil {
		return err
	}

	return nil
}

func (ds *DocStore) DeltaLogSize(ctx context.Context, docID string) (int64, error) {
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	size, err := ds.Redis.ZCard(ctx, eventsZSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("DeltaLogSize(%q): %w", docID, err)
	}
	return size, nil
}

func (ds *DocStore) GetDeltaLog(ctx context.Context, docID string, lowScore, highScore int64) ([]redis.Z, error) {
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	max := fmt.Sprintf("%d", highScore)
	if highScore == -1 {
		max = "+inf"
	}

	// ZRangeByScore gets all elements in the sorted set within the specified score range
	events, err := ds.Redis.ZRangeByScoreWithScores(ctx, eventsZSetKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", lowScore),
		Max: max,
	}).Result()

	if err != nil {
		return nil, err
	}

	log.Infof("GetDeltaLog(%q, %d, %d) -> %d events", docID, lowScore, highScore, len(events))

	return events, nil
}

func (ds *DocStore) DeleteDeltaLogItem(ctx context.Context, docID string, score int64) error {
	eventsZSetKey := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	members, err := ds.Redis.ZRangeByScore(ctx, eventsZSetKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", score),
		Max: fmt.Sprintf("%d", score),
	}).Result()
	if err != nil {
		return fmt.Errorf("error finding items by score: %v", err)
	}

	if len(members) > 0 {
		_, err = ds.Redis.ZRemRangeByScore(ctx, eventsZSetKey, fmt.Sprintf("%d", score), fmt.Sprintf("%d", score)).Result()
		if err != nil {
			return fmt.Errorf("error deleting original items: %v", err)
		}

		_, err = ds.Redis.ZAdd(ctx, eventsZSetKey, redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("deleted-%d", score),
		}).Result()
		if err != nil {
			return fmt.Errorf("error adding 'deleted' item: %v", err)
		}
	} else {
		log.Info("No items found with the specified score to delete.")
	}

	return nil
}

func (ds *DocStore) SizeDeltaLog(ctx context.Context, docID string) (int64, error) {
	key := fmt.Sprintf(constants.DocEventsKeyFormat, docID)

	size, err := ds.Redis.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return size, nil
}

type ErrorPrefixNotFound struct {
	Prefix string
}

func (e ErrorPrefixNotFound) Error() string {
	return fmt.Sprintf("no prefix: %s", e.Prefix)
}

func (ds *DocStore) ListFirst(prefix string) (string, error) {
	input := &awsS3.ListObjectsV2Input{
		Bucket:  aws.String(ds.S3.Bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(1),
	}

	result, err := ds.S3.Client.ListObjectsV2(input)
	if err != nil {
		log.Errorf("failed to list objects: %s", err)
		return "", err
	}

	if len(result.Contents) == 0 {
		return "", ErrorPrefixNotFound{Prefix: prefix}
	}

	return *result.Contents[0].Key, nil
}

func (ds *DocStore) GetFirst(prefix string) ([]byte, error) {
	key, err := ds.ListFirst(prefix)
	if err != nil {
		return []byte(""), err
	}

	bytes, err := ds.S3.GetObject(ds.S3.Bucket, key)
	if err != nil {
		return []byte(""), err
	}

	return bytes, nil
}

func invertSeq(seq int64) string {
	invertedSeq := constants.MaxSeqValue - seq
	return fmt.Sprintf("%016d", invertedSeq)
}

func DocS3Path(docID string) string {
	return fmt.Sprintf(constants.DocumentSnapshotPrefix, constants.S3Prefix, docID)
}

func SeqToS3Path(docID string, seq int64) string {
	return path.Join(DocS3Path(docID), invertSeq(seq))
}

func revertInvertedSeq(invSeq string) (int64, error) {
	invertedSeq, err := strconv.ParseInt(invSeq, 10, 64)
	if err != nil {
		return 0, err
	}
	return constants.MaxSeqValue - invertedSeq, nil
}

func ExtractSeqFromS3Path(s3Path string) (int64, error) {
	parts := strings.Split(s3Path, "/")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid S3 path: %s", s3Path)
	}
	invSeq := parts[3]

	return revertInvertedSeq(invSeq)
}

func (ds *DocStore) GetLastS3Seq(docID string) (int64, error) {
	snapPath := DocS3Path(docID)
	snapKey, err := ds.ListFirst(snapPath)
	if err != nil {
		if errors.As(err, &ErrorPrefixNotFound{}) {
			return 0, nil
		} else {
			return 0, err
		}
	}

	seq, err := ExtractSeqFromS3Path(snapKey)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

func (ds *DocStore) GetLastSnapshot(docID string) (int64, *rogueV3.Rogue, error) {
	snapPath := DocS3Path(docID)
	snapKey, err := ds.ListFirst(snapPath)
	if err != nil {
		if errors.As(err, &ErrorPrefixNotFound{}) {
			// Handle the prefix not found error
			return 0, rogueV3.NewRogueForQuill("nf"), nil
		} else {
			return 0, nil, err
		}
	}

	snap, err := ds.S3.GetObject(ds.S3.Bucket, snapKey)
	if err != nil {
		return 0, nil, err
	}

	seq, err := ExtractSeqFromS3Path(snapKey)
	if err != nil {
		return 0, nil, err
	}

	var doc rogueV3.Rogue
	err = json.Unmarshal(snap, &doc)
	if err != nil {
		return 0, nil, fmt.Errorf("error unmarshalling snapshot: %w", err)
	}

	return seq, &doc, nil
}

func (ds *DocStore) GetCurrentDoc(ctx context.Context, docID string) (int64, *rogueV3.Rogue, error) {
	seq, doc, err := ds.GetLastSnapshot(docID)
	if err != nil {
		return 0, nil, fmt.Errorf("ds.GetLastSnapshot(%s): %w", docID, err)
	}

	ops, err := ds.GetDeltaLog(ctx, docID, 0, -1)
	if err != nil {
		return 0, nil, fmt.Errorf("ds.GetDeltaLog(ctx, %s, 0, -1): %w", docID, err)
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

		var msg rogueV3.Message
		err = json.Unmarshal([]byte(op), &msg)
		if err != nil {
			log.Errorf("[doc %s] error unmarshalling op: %s", docID, err)
			return 0, nil, err
		}

		_, err := doc.MergeOp(msg.Op)
		if err != nil {
			log.Errorf("[doc %s] Error merging op %#v: %s", docID, msg.Op, err)
			// return 0, nil, err
		}
		maxSeq = int64(z.Score)
	}

	return maxSeq, doc, nil
}

func (ds *DocStore) SnapshotDoc(ctx context.Context, docID string) error {
	seq, doc, err := ds.GetCurrentDoc(ctx, docID)
	if err != nil {
		return err
	}

	err = ds.SaveDocToS3(ctx, docID, seq, doc)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	err = ds.CleanupDeltaLog(ctx, docID, seq)
	if err != nil {
		return err
	}

	return nil
}

// SaveDocToS3 saves the doc to S3, it shouldn't normally be called directly. But is used directly in tests.
func (ds *DocStore) SaveDocToS3(ctx context.Context, docID string, seq int64, doc *rogueV3.Rogue) error {
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshalling doc: %s", err)
	}

	s3path := SeqToS3Path(docID, seq)
	log.Info("Saving doc to S3", "s3path", s3path)
	err = ds.S3.PutObject(ds.S3.Bucket, s3path, "text/plain", docBytes)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	return nil
}

func (ds *DocStore) DuplicateDoc(ctx context.Context, docID, newDocID string, addressString *string) error {
	seq, curDoc, err := ds.GetCurrentDoc(ctx, docID)
	if err != nil {
		return err
	}

	var address *rogueV3.ContentAddress = nil
	if addressString != nil {
		address = &rogueV3.ContentAddress{}
		err = json.Unmarshal([]byte(*addressString), address)
		if err != nil {
			return stackerr.Wrap(err)
		}
	}

	dupDocRogue, err := curDoc.Compact(address)
	if err != nil {
		return err
	}

	err = ds.SaveDocToS3(ctx, newDocID, seq, dupDocRogue)
	if err != nil {
		return err
	}

	err = ds.CleanupDeltaLog(ctx, newDocID, seq)
	if err != nil {
		return err
	}

	return nil
}
