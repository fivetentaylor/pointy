package dynamo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/teamreviso/code/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

const DefaultTimelineChain = "main"

var TimelinePrefix = "time#"

type TimelineEvent struct {
	// PK
	DocID string `json:"docID"`

	// SK
	EventID        string `json:"eventID"`
	DefaultEventID string `json:"-"`

	// SK1
	ReplyToID string `json:"chain"`
	CreatedAt int64  `json:"createdAt"`

	// Attributes
	UserID   string `json:"userID"`
	AuthorID string `json:"authorID"`

	// Payload
	Event *models.TimelineEventPayload `json:"event"`
}

type dTimelineEvent struct {
	// doc#{uuid}
	PK string `dynamodbav:"PK" json:"PK"`
	// time#{uuid}
	SK string `dynamodbav:"SK" json:"SK"`
	// time#main@{createdAt} or time#reply#{uuid}@{createdAt}
	SK1 string `dynamodbav:"SK1" json:"SK1"`

	UserID   string `dynamodbav:"userID" json:"userID"`
	AuthorID string `dynamodbav:"authorID" json:"authorID"`

	Event string `dynamodbav:"event" json:"event"`
}

func (db *DB) CreateTimelineEvent(event *TimelineEvent) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	if event.DocID == "" {
		return fmt.Errorf("missing docID")
	}
	if event.EventID != "" {
		return fmt.Errorf("already has eventID")
	}

	if event.DefaultEventID != "" {
		event.EventID = event.DefaultEventID
	} else {
		event.EventID = uuid.NewString()
	}
	event.CreatedAt = time.Now().UnixNano()

	av, err := dynamodbattribute.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      av.M,
	}

	_, err = db.Client.PutItem(putInput)
	if err != nil {
		return fmt.Errorf("CreateTimelineEvent(%v): %w", event, err)
	}

	return nil
}

func (db *DB) UpdateTimelineEvent(event *TimelineEvent) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	if event.EventID == "" {
		return fmt.Errorf("missing eventID")
	}

	av, err := dynamodbattribute.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      av.M,
	}

	_, err = db.Client.PutItem(putInput)
	if err != nil {
		return fmt.Errorf("UpdateTimelineEvent(%v): %w", event, err)
	}

	return nil
}

func (db *DB) GetTimelineEvent(docId, eventID string) (*TimelineEvent, error) {
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("doc#%s", docId)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("time#%s", eventID)),
			},
		},
	}

	result, err := db.Client.GetItem(getInput)
	if err != nil {
		return nil, fmt.Errorf("GetTimelineEvent(%q, %q): %w", docId, eventID, err)
	}

	if result.Item == nil {
		return nil, nil
	}

	event := &TimelineEvent{}
	err = dynamodbattribute.UnmarshalMap(result.Item, event)
	if err != nil {
		return nil, fmt.Errorf("GetTimelineEvent(%q, %q): %w", docId, eventID, err)
	}
	return event, nil
}

func (db *DB) DeleteTimelineEvent(docId, eventID string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("doc#%s", docId)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("time#%s", eventID)),
			},
		},
	}

	_, err := db.Client.DeleteItem(deleteInput)
	if err != nil {
		return fmt.Errorf("DeleteTimelineEvent(%q, %q): %w", docId, eventID, err)
	}

	return nil
}

func (db *DB) GetDocumentTimeline(docId string) ([]*TimelineEvent, error) {
	return db.getTimelineEvents(docId, DefaultTimelineChain)
}

func (db *DB) GetDocumentTimelineReplies(docId string, eventID string) ([]*TimelineEvent, error) {
	return db.getTimelineEvents(docId, fmt.Sprintf("reply#%s", eventID))
}

func (db *DB) GetLastUserUpdate(docId, userID string) (*TimelineEvent, error) {
	events, err := db.getTimelineEvents(docId, DefaultTimelineChain)
	if err != nil {
		return nil, fmt.Errorf("GetLastUserUpdate(%q, %q): %w", docId, userID, err)
	}

	for i := len(events) - 1; i >= 0; i-- {
		if events[i].UserID != userID {
			continue
		}

		if events[i].Event.GetUpdate() == nil {
			continue
		}

		return events[i], nil
	}

	return nil, nil
}

func (db *DB) GetLastCompletedUserUpdate(docId, userID string) (*TimelineEvent, error) {
	events, err := db.getTimelineEvents(docId, DefaultTimelineChain)
	if err != nil {
		return nil, fmt.Errorf("GetLastCompletedUserUpdate(%q, %q): %w", docId, userID, err)
	}

	for i := len(events) - 1; i >= 0; i-- {
		if events[i].UserID != userID {
			continue
		}

		if events[i].Event.GetUpdate() == nil {
			continue
		}

		if events[i].Event.GetUpdate().State != models.UpdateState_COMPLETE_STATE {
			continue
		}

		return events[i], nil
	}

	return nil, nil
}

func (e *TimelineEvent) GetKey() map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(e.DocID)},
		"SK": {S: aws.String(fmt.Sprintf("%s%s", TimelinePrefix, e.EventID))},
	}

	return key
}

func (e *TimelineEvent) CreatedAtTimestamp() time.Time {
	return time.UnixMicro(e.CreatedAt / 1000)
}

func (e *TimelineEvent) UpdateState() models.UpdateState {
	if e.Event.GetUpdate() == nil {
		return models.UpdateState_UNKNOWN_STATE
	}

	return e.Event.GetUpdate().State
}

func (e *TimelineEvent) MarshalJSON() ([]byte, error) {
	marshalItem, err := e.toDynamo()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal internal message: %s", err)
	}

	bts, err := json.Marshal(marshalItem)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal internal message: %s", err)
	}
	return bts, nil
}

func (e *TimelineEvent) UnmarshalJSON(b []byte) error {
	tmp := dTimelineEvent{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := e.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (e TimelineEvent) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	marshalItem, err := e.toDynamo()
	if err != nil {
		return err
	}

	ma, err := dynamodbattribute.MarshalMap(marshalItem)
	if err != nil {
		return err
	}
	av.M = ma

	return nil
}

func (e *TimelineEvent) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dTimelineEvent{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := e.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (e *TimelineEvent) toDynamo() (dTimelineEvent, error) {
	chain := DefaultTimelineChain

	if e.ReplyToID != "" {
		chain = fmt.Sprintf("reply#%s", e.ReplyToID)
	}

	dte := dTimelineEvent{
		PK:       fmt.Sprintf("doc#%s", e.DocID),
		SK:       fmt.Sprintf("%s%s", TimelinePrefix, e.EventID),
		SK1:      fmt.Sprintf("%s%s@%d", TimelinePrefix, chain, e.CreatedAt),
		UserID:   e.UserID,
		AuthorID: e.AuthorID,
	}

	eventBts, err := protojson.Marshal(e.Event)
	if err != nil {
		return dte, fmt.Errorf("failed to marshal event: %s", err)
	}
	dte.Event = string(eventBts)

	return dte, nil
}

func (e *TimelineEvent) fromDynamo(tmp dTimelineEvent) error {
	var err error

	e.UserID = tmp.UserID
	e.AuthorID = tmp.AuthorID

	e.DocID = strings.TrimPrefix(tmp.PK, DocPrefix)
	e.EventID = strings.TrimPrefix(tmp.SK, TimelinePrefix)

	sk1Parts := strings.Split(
		strings.TrimPrefix(tmp.SK1, TimelinePrefix),
		"@",
	)
	if len(sk1Parts) != 2 {
		return fmt.Errorf("invalid sk1: %s", tmp.SK1)
	}
	if sk1Parts[0] != DefaultTimelineChain {
		e.ReplyToID = strings.TrimPrefix(sk1Parts[0], "reply#")
	}

	e.CreatedAt, err = strconv.ParseInt(sk1Parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse createdAt: %s", err)
	}

	e.Event = &models.TimelineEventPayload{}
	err = protoJsonUnmarshaler.Unmarshal([]byte(tmp.Event), e.Event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal timeline event: %s/%s %s", tmp.PK, tmp.SK, err)
	}

	return nil
}

func (db *DB) getTimelineEvents(docId, chain string) ([]*TimelineEvent, error) {
	query := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK1, :prefix))"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("doc#%s", docId))},
			":prefix": {S: aws.String(fmt.Sprintf("%s%s", TimelinePrefix, chain))},
		},
		ScanIndexForward: aws.Bool(true),
	}

	results, err := db.Client.Query(query)
	if err != nil {
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, nil
	}

	events := make([]*TimelineEvent, len(results.Items))
	for i, item := range results.Items {
		msg := TimelineEvent{}
		err := dynamodbattribute.UnmarshalMap(item, &msg)
		if err != nil {
			log.Errorf("failed to unmarshal message map: %s", err)
			return nil, err
		}
		events[i] = &msg
	}

	return events, nil
}
