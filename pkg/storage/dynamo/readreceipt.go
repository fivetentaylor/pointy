package dynamo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
)

var ErrReadReceiptNotFound = fmt.Errorf("read receipt not found")

var readReceiptPrefix = "rr#"

type ReadReceipt struct {
	// PK
	UserID string `json:"userID"`
	// SK
	MessageID string `json:"messageID"`
	// SK1
	Read        bool   `json:"read"`
	DocID       string `json:"docID"`
	ChannelID   string `json:"channelID"`
	ContainerID string `json:"containerID"`
	// SK2
	Mentioned bool `json:"mentioned"`

	// Attributes
	CreatedAt int64 `json:"createdAt"`
}

type dReadReceipt struct {
	// rr#{UserID}
	PK string `dynamodbav:"PK"`
	// {MessageID}
	SK string `dynamodbav:"SK"`
	// rr#{Read}#{DocID}#{ChannelID}#{MessageID}
	SK1 string `dynamodbav:"SK1"`

	CreatedAt int64 `dynamodbav:"createdAt"`
}

func (item ReadReceipt) Key() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", readReceiptPrefix, item.UserID))},
		"SK": {S: aws.String(item.MessageID)},
	}
}

func (item ReadReceipt) SK1() string {
	return fmt.Sprintf("%s%t#%s#%s#%s", readReceiptPrefix, item.Read, item.DocID, item.ChannelID, item.ContainerID)
}

func (item ReadReceipt) SK2() string {
	return fmt.Sprintf("%s%t#%t#%s#%s#%s", readReceiptPrefix, item.Mentioned, item.Read, item.DocID, item.ChannelID, item.ContainerID)
}

func (item *ReadReceipt) HydrateFromSK1(sk1 string) error {
	var err error
	parts := strings.Split(sk1, "#")

	item.Read, err = strconv.ParseBool(parts[1])
	if err != nil {
		return fmt.Errorf("HydrateFromSK1(%s): %w", sk1, err)
	}
	item.DocID = parts[2]
	item.ChannelID = parts[3]
	item.ContainerID = strings.Join(parts[4:], "#")
	return nil
}

func (item ReadReceipt) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	key := item.Key()
	marshalItem := dReadReceipt{
		PK:        *key["PK"].S,
		SK:        *key["SK"].S,
		SK1:       item.SK1(),
		CreatedAt: item.CreatedAt,
	}
	m, err := dynamodbattribute.MarshalMap(marshalItem)
	if err != nil {
		return err
	}
	av.M = m
	return nil
}

func (item *ReadReceipt) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dReadReceipt{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return err
	}

	item.UserID = strings.TrimPrefix(tmp.PK, readReceiptPrefix)
	err := item.HydrateFromSK1(tmp.SK1)
	if err != nil {
		return fmt.Errorf("UnmarshalDynamoDBAttributeValue: %w", err)
	}

	item.MessageID = tmp.SK
	item.CreatedAt = tmp.CreatedAt

	return nil
}

func (db *DB) FindOrCreateReadReceipt(userID, docID, channelID, containerID, messageID string, mentioned bool) (*ReadReceipt, error) {
	rr := &ReadReceipt{
		UserID:      userID,
		DocID:       docID,
		ContainerID: containerID,
		ChannelID:   channelID,
		MessageID:   messageID,
		Read:        false,
		Mentioned:   mentioned,
	}

	key := rr.Key()

	// Attempt to find the existing item
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       key,
	}

	result, err := db.Client.GetItem(getItemInput)
	if err != nil {
		return nil, fmt.Errorf("FindOrCreateReadReceipt: failed to get item: %s", err)
	}

	// If item exists, unmarshal and return it
	if result.Item != nil {
		var existingRR ReadReceipt
		err = dynamodbattribute.UnmarshalMap(result.Item, &existingRR)
		if err != nil {
			return nil, fmt.Errorf("FindOrCreateReadReceipt: failed to unmarshal existing item: %s", err)
		}
		return &existingRR, nil
	}

	// Item does not exist, create it
	rr.CreatedAt = time.Now().UnixNano()
	av, err := dynamodbattribute.MarshalMap(rr)
	if err != nil {
		return nil, fmt.Errorf("FindOrCreateReadReceipt: failed to marshal rr: %s", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(db.TableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err = db.Client.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("FindOrCreateReadReceipt: failed to put item: %s", err)
	}

	return rr, nil
}

func (db *DB) GetReadReceipt(userID, messageID string) (*ReadReceipt, error) {
	rr := &ReadReceipt{
		UserID:    userID,
		MessageID: messageID,
	}

	key := rr.Key()

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       key,
	}

	result, err := db.Client.GetItem(getItemInput)
	if err != nil {
		return nil, fmt.Errorf("GetReadReceipt(%q, %q): failed to get item: %w", userID, messageID, err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("GetReadReceipt(%q, %q): %w", userID, messageID, ErrReadReceiptNotFound)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, rr)
	if err != nil {
		return nil, fmt.Errorf("GetReadReceipt(%q, %q): failed to get item: %w", userID, messageID, err)
	}
	return rr, nil
}

func (db *DB) GetReadReceiptsForUser(userID string) ([]ReadReceipt, error) {
	pk := fmt.Sprintf("%s%s", readReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK"),
		IndexName:              aws.String("SK1"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
		},
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unreads: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, fmt.Errorf("unread not found")
	}

	if len(results.Items) > 1 {
		log.Errorf("multiple unread found")
	}

	unreads := make([]ReadReceipt, len(results.Items))
	for i, item := range results.Items {
		urd := ReadReceipt{}
		err := dynamodbattribute.UnmarshalMap(item, &urd)
		if err != nil {
			log.Errorf("failed to unmarshal message map: %s", err)
			return nil, err
		}
		unreads[i] = urd
	}

	return unreads, nil

}

func (db *DB) GetUnreadCountForContainer(userID, docID, channelID, containerID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%s#%s#%s", readReceiptPrefix, false, docID, channelID, containerID)
	return db.getUnreadCount(userID, skPrefix)
}

func (db *DB) GetUnreadMentionCountForContainer(userID, docID, channelID, containerID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%t#%s#%s#%s", readReceiptPrefix, true, false, docID, channelID, containerID)
	return db.getUnreadMentionCount(userID, skPrefix)
}

func (db *DB) GetUnreadCountForChannel(userID, docID, channelID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%s#%s", readReceiptPrefix, false, docID, channelID)
	return db.getUnreadCount(userID, skPrefix)
}

func (db *DB) GetUnreadMentionCountForChannel(userID, docID, channelID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%t#%s#%s", readReceiptPrefix, true, false, docID, channelID)
	return db.getUnreadMentionCount(userID, skPrefix)
}

func (db *DB) GetUnreadCountForDocument(userID, docID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%s", readReceiptPrefix, false, docID)
	return db.getUnreadCount(userID, skPrefix)
}

func (db *DB) GetUnreadMentionCountForDocument(userID, docID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%t#%s", readReceiptPrefix, true, false, docID)
	return db.getUnreadMentionCount(userID, skPrefix)
}

func (db *DB) GetUnreadCountForUser(userID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t", readReceiptPrefix, false)
	return db.getUnreadCount(userID, skPrefix)
}
func (db *DB) GetUnreadMentionCountForUser(userID string) (int64, error) {
	skPrefix := fmt.Sprintf("%s%t#%t", readReceiptPrefix, true, false)
	return db.getUnreadMentionCount(userID, skPrefix)
}

func (db *DB) MarkReadReceiptRead(userID, docID, channelID, containerID, messageID string) (*ReadReceipt, error) {
	rr, err := db.UpdateReadReceiptRead(userID, docID, channelID, containerID, messageID, true)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			// if they're trying to mark a read receipt that doesn't exist yet, we don't really care about if they were mentioned or not
			rr, err = db.FindOrCreateReadReceipt(userID, docID, channelID, containerID, messageID, false)
			if err != nil {
				return rr, err
			}
			return db.UpdateReadReceiptRead(userID, docID, channelID, containerID, messageID, true) // TODO this should be in the create
		}
		return nil, err
	}

	return rr, nil
}

func (db *DB) MarkReadReceiptUnread(userID, docID, channelID, containerID, messageID string, mentioned bool) (*ReadReceipt, error) {
	rr, err := db.UpdateReadReceiptRead(userID, docID, channelID, containerID, messageID, false)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return db.FindOrCreateReadReceipt(userID, docID, channelID, containerID, messageID, mentioned)
		}
		return nil, err
	}

	return rr, nil
}

func (db *DB) UpdateReadReceiptRead(userID, docID, channelID, containerID, messageID string, read bool) (*ReadReceipt, error) {
	rr := &ReadReceipt{
		UserID:      userID,
		DocID:       docID,
		ChannelID:   channelID,
		ContainerID: containerID,
		MessageID:   messageID,
		Read:        read,
	}

	av, err := dynamodbattribute.MarshalMap(rr)
	if err != nil {
		return nil, err
	}

	updateExpression := "SET #SK1 = :SK1, #createdAt = :createdAt"
	exprAttrNames := map[string]*string{
		"#SK1":       aws.String("SK1"),
		"#createdAt": aws.String("createdAt"),
	}
	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":SK1":       av["SK1"],
		":createdAt": av["createdAt"],
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.TableName),
		Key:                       rr.Key(),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ConditionExpression:       aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
		ReturnValues:              aws.String("NONE"),
	}

	// Perform the update operation
	_, err = db.Client.UpdateItem(updateItemInput)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (db *DB) getReadReceipts(userID, skPrefix string) ([]ReadReceipt, error) {
	pk := fmt.Sprintf("%s%s", readReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK1Index"),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK1, :SKPrefix)"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":       {S: aws.String(pk)},
			":SKPrefix": {S: aws.String(skPrefix)},
		},
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unreads: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, fmt.Errorf("unread not found")
	}

	if len(results.Items) > 1 {
		log.Errorf("multiple unread found")
	}

	unreads := make([]ReadReceipt, len(results.Items))
	for i, item := range results.Items {
		urd := ReadReceipt{}
		err := dynamodbattribute.UnmarshalMap(item, &urd)
		if err != nil {
			log.Errorf("failed to unmarshal message map: %s", err)
			return nil, err
		}
		unreads[i] = urd
	}

	return unreads, nil
}

func (db *DB) getUnreadCount(userID, skPrefix string) (int64, error) {
	pk := fmt.Sprintf("%s%s", readReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK1Index"),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK1, :SKPrefix)"),
		Select:                 aws.String("COUNT"), // Change to COUNT to get only the number of items
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":       {S: aws.String(pk)},
			":SKPrefix": {S: aws.String(skPrefix)},
		},
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unreads count: %s", err)
		return 0, err
	}

	return *results.Count, nil
}

func (db *DB) getUnreadMentionCount(userID, skPrefix string) (int64, error) {
	pk := fmt.Sprintf("%s%s", readReceiptPrefix, userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK2Index"),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK2, :SKPrefix)"),
		Select:                 aws.String("COUNT"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":       {S: aws.String(pk)},
			":SKPrefix": {S: aws.String(skPrefix)},
		},
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query unreads count: %s", err)
		return 0, err
	}

	return *results.Count, nil
}
