package dynamo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
)

var ErrSentReceiptNotFound = fmt.Errorf("sent receipt not found")

var sentReceiptPrefix = "sr#"

type SentReceipt struct {
	// PK
	UserID string `json:"userID"`
	// SK
	MessageID string `json:"messageID"`
	// SK1
	SentAt int64 `json:"sentAt"`

	ContainerID string `json:"containerID"`
}

type dSentReceipt struct {
	// sr#{UserID}
	PK string `dynamodbav:"PK"`
	// {MessageID}
	SK string `dynamodbav:"SK"`
	// {sentAt}
	SK1 string `dynamodbav:"SK1"`

	ContainerID string `dynamodbav:"containerID"`
}

func (item SentReceipt) Key() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", sentReceiptPrefix, item.UserID))},
		"SK": {S: aws.String(item.MessageID)},
	}
}

func (item SentReceipt) SentAtTimestamp() time.Time {
	return time.UnixMicro(item.SentAt / 1000)
}

func (item SentReceipt) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	key := item.Key()
	marshalItem := dSentReceipt{
		PK:          *key["PK"].S,
		SK:          *key["SK"].S,
		SK1:         fmt.Sprintf("%d", item.SentAt),
		ContainerID: item.ContainerID,
	}
	m, err := dynamodbattribute.MarshalMap(marshalItem)
	if err != nil {
		return err
	}
	av.M = m
	return nil
}

func (item *SentReceipt) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dSentReceipt{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return err
	}

	var err error
	item.UserID = strings.TrimPrefix(tmp.PK, sentReceiptPrefix)
	item.MessageID = tmp.SK
	item.ContainerID = tmp.ContainerID
	item.SentAt, err = strconv.ParseInt(tmp.SK1, 10, 64)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateSentReceipt(userID, messageID, containerID string) (*SentReceipt, error) {
	sr := &SentReceipt{
		UserID:      userID,
		MessageID:   messageID,
		ContainerID: containerID,
		SentAt:      0,
	}

	av, err := dynamodbattribute.MarshalMap(sr)

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(db.TableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err = db.Client.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("CreateSentReceipt: failed to put item: %s", err)
	}

	return sr, nil
}

func (db *DB) GetSentReceipt(userID, messageID string) (*SentReceipt, error) {
	sr := &SentReceipt{
		UserID:    userID,
		MessageID: messageID,
	}

	key := sr.Key()

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       key,
	}

	result, err := db.Client.GetItem(getItemInput)
	if err != nil {
		return nil, fmt.Errorf("GetSentReceipt(%q, %q): failed to get item: %w", userID, messageID, err)
	}

	if result.Item == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, sr)
	if err != nil {
		return nil, fmt.Errorf("GetSentReceipt(%q, %q): failed to get item: %w", userID, messageID, err)
	}
	return sr, nil
}

func (db *DB) GetUnsentSentReceiptsForUser(userID string) ([]SentReceipt, error) {
	pk := fmt.Sprintf("%s%s", sentReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
		},
		Limit: aws.Int64(10),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unsents: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return []SentReceipt{}, nil
	}

	unsents := make([]SentReceipt, 0, len(results.Items))
	for _, item := range results.Items {
		urd := SentReceipt{}
		err := dynamodbattribute.UnmarshalMap(item, &urd)
		if err != nil {
			log.Errorf("failed to unmarshal message map: %s", err)
			return nil, err
		}
		if urd.SentAt != 0 {
			continue
		}

		unsents = append(unsents, urd)
	}

	return unsents, nil
}

func (db *DB) UpdateSentReceiptSent(userID, messageID string, sent bool) (*SentReceipt, error) {
	sr := &SentReceipt{
		UserID:    userID,
		MessageID: messageID,
		SentAt:    0,
	}

	if sent {
		sr.SentAt = time.Now().UnixNano()
	}

	av, err := dynamodbattribute.MarshalMap(sr)
	if err != nil {
		return nil, err
	}

	updateExpression := "SET #SK1 = :SK1"
	exprAttrNames := map[string]*string{
		"#SK1": aws.String("SK1"),
	}
	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":SK1": av["SK1"],
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.TableName),
		Key:                       sr.Key(),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ConditionExpression:       aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
		ReturnValues:              aws.String("NONE"),
	}

	log.Infof("UpdateSentReceiptSent: %v", updateItemInput)

	// Perform the update operation
	_, err = db.Client.UpdateItem(updateItemInput)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func (db *DB) DeleteSentReceipt(userID, messageID string) error {
	sr := &SentReceipt{
		UserID:    userID,
		MessageID: messageID,
	}

	key := sr.Key()
	_, err := db.Client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       key,
	})
	return err
}

func (db *DB) GetLastSentReceipt(userID string) (*SentReceipt, error) {
	pk := fmt.Sprintf("%s%s", sentReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
		},
		Limit:            aws.Int64(1),
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unsents: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, nil
	}

	urd := SentReceipt{}
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &urd)
	if err != nil {
		log.Errorf("failed to unmarshal message map: %s", err)
		return nil, err
	}

	return &urd, nil
}
