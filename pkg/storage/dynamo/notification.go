package dynamo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fivetentaylor/pointy/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

// Access Patterns:
// - Get unread notifications for user (order by createdAt)
// - Get read notifications for user (order by createdAt)
// - Get read notifications for user & document (order by createdAt)
// - Get unread notifications for user & document (order by createdAt)
// - Get a count of unread notifications for user
// - Get a count of unread notifications for user & document

var ErrNotificationNotFound = errors.New("notification not found")

var NotificationPrefix = "notif#"

type Notification struct {
	// PK
	UserID string `json:"userID"`
	// SK
	ID string `json:"id"`
	// SK1/2
	DocID     string `json:"docID"`
	Read      bool   `json:"read"`
	CreatedAt int64  `json:"createdAt"`

	// Attributes
	Payload *models.NotificationPayload `json:"payload"`
}

func (n Notification) toDynamo() (dNotification, error) {
	dn := dNotification{
		PK:  fmt.Sprintf("%s%s", NotificationPrefix, n.UserID),
		SK:  n.ID,
		SK1: fmt.Sprintf("%t#%d", n.Read, n.CreatedAt),
		SK2: fmt.Sprintf("%s#%t#%d", n.DocID, n.Read, n.CreatedAt),
	}

	bts, err := protojson.Marshal(n.Payload)
	if err != nil {
		return dn, fmt.Errorf("failed to marshal payload: %s", err)
	}

	dn.Payload = string(bts)
	return dn, nil
}

func (n *Notification) fromDynamo(dn dNotification) error {
	var err error
	n.UserID = strings.TrimPrefix(dn.PK, NotificationPrefix)
	n.ID = dn.SK

	sk2Parts := strings.Split(dn.SK2, "#")
	n.DocID = sk2Parts[0]
	n.Read, err = strconv.ParseBool(sk2Parts[1])
	if err != nil {
		return err
	}
	n.CreatedAt, err = strconv.ParseInt(sk2Parts[2], 10, 64)
	if err != nil {
		return err
	}

	p := &models.NotificationPayload{}
	if err := protojson.Unmarshal([]byte(dn.Payload), p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %s", err)
	}
	n.Payload = p

	return nil
}

type dNotification struct {
	// PK: notif#userID
	PK string `dynamodbav:"PK" json:"PK"`
	// SK: messageID
	SK string `dynamodbav:"SK" json:"SK"`
	// SK1: read#createdAt
	SK1 string `dynamodbav:"SK1" json:"SK1"`
	// SK2: docID#read#createdAt
	SK2 string `dynamodbav:"SK2" json:"SK2"`

	Payload string `dynamodbav:"payload" json:"payload"`
}

func (n Notification) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	marshalItem, err := n.toDynamo()
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

func (n *Notification) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dNotification{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := n.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (db *DB) GetNotification(userID, id string) (*Notification, error) {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)
	sk := id

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(pk)},
			"SK": {S: aws.String(sk)},
		},
	}

	result, err := db.Client.GetItem(getInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %s", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("%s-%s: %w", pk, sk, ErrNotificationNotFound)
	}

	n := Notification{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &n); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %s", err)
	}

	return &n, nil
}

func (db *DB) UpsertNotification(n *Notification) error {
	if n.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	if n.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if n.DocID == "" {
		return fmt.Errorf("docID cannot be empty")
	}
	n.CreatedAt = time.Now().UnixNano()

	av, err := dynamodbattribute.MarshalMap(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %s", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      av,
	}

	_, err = db.Client.PutItem(input)
	if err != nil {
		return fmt.Errorf("UpsertNotification(%+v) (%+v): %w", n, input, err)
	}

	return nil
}

func (db *DB) GetNotificationsForUser(userID string, read bool, params PaginationParams) ([]Notification, map[string]*dynamodb.AttributeValue, error) {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)
	sk := fmt.Sprintf("%t", read)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK1, :SK)"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
			":SK": {S: aws.String(sk)},
		},
		Limit:             aws.Int64(params.Limit),
		ExclusiveStartKey: params.ExclusiveStartKey,
		ScanIndexForward:  aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return nil, nil, fmt.Errorf("GetNotificationsByUser(%q, %+v): %s", userID, params, err)
	}

	notifications := make([]Notification, len(results.Items))
	for i, item := range results.Items {
		n := Notification{}
		err := dynamodbattribute.UnmarshalMap(item, &n)
		if err != nil {
			return nil, nil, fmt.Errorf("GetNotificationsByUser(%q, %+v): %s", userID, params, err)
		}
		notifications[i] = n
	}

	return notifications, results.LastEvaluatedKey, nil
}

func (db *DB) GetNotificationsForUserDocument(userID, docID string, read bool, params PaginationParams) ([]*Notification, map[string]*dynamodb.AttributeValue, error) {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)
	sk2 := fmt.Sprintf("%s#%t", docID, read)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK2, :SK)"),
		IndexName:              aws.String("SK2Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(pk)},
			":SK": {S: aws.String(sk2)},
		},
		Limit:             aws.Int64(params.Limit),
		ExclusiveStartKey: params.ExclusiveStartKey,
		ScanIndexForward:  aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return nil, nil, fmt.Errorf("GetNotificationsForUserDocument(%q, %q, %+v): %s", userID, docID, params, err)
	}

	notifications := make([]*Notification, len(results.Items))
	for i, item := range results.Items {
		n := &Notification{}
		err := dynamodbattribute.UnmarshalMap(item, n)
		if err != nil {
			return nil, nil, fmt.Errorf("GetNotificationsForUserDocument(%q, %q, %+v): %s", userID, docID, params, err)
		}
		notifications[i] = n
	}

	return notifications, results.LastEvaluatedKey, nil
}

func (db *DB) GetNotificationCountForUser(userID string, read bool) (int64, error) {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)
	sk1 := fmt.Sprintf("%t", read)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK1Index"),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK1, :SK1)"),
		Select:                 aws.String("COUNT"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":  {S: aws.String(pk)},
			":SK1": {S: aws.String(sk1)},
		},
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return 0, fmt.Errorf("GetUnreadNotificationCountForUser(%q): %s", userID, err)
	}

	return *results.Count, nil
}

func (db *DB) GetNotificationCountForUserAndDocument(userID, docID string, read bool) (int64, error) {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)
	sk2 := fmt.Sprintf("%s#%t", docID, read)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK2Index"),
		KeyConditionExpression: aws.String("PK = :PK AND begins_with(SK2, :SK2)"),
		Select:                 aws.String("COUNT"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":  {S: aws.String(pk)},
			":SK2": {S: aws.String(sk2)},
		},
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return 0, fmt.Errorf("GetUnreadNotificationCountForUser(%q): %s", userID, err)
	}

	return *results.Count, nil
}

func (db *DB) MarkNotifications(userID, docID string, ids []string, newState bool) error {
	pk := fmt.Sprintf("%s%s", NotificationPrefix, userID)

	transactItems := []*dynamodb.TransactWriteItem{}

	for _, id := range ids {
		notif, err := db.GetNotification(userID, id)
		if err != nil {
			return fmt.Errorf("failed to get notification: %s", err)
		}

		sk := id
		updateExpr := "SET SK1 = :SK1, SK2 = :SK2"
		expressionAttributeValues := map[string]*dynamodb.AttributeValue{
			":SK1": {S: aws.String(fmt.Sprintf("%t#%d", newState, notif.CreatedAt))},
			":SK2": {S: aws.String(fmt.Sprintf("%s#%t#%d", docID, newState, notif.CreatedAt))},
		}

		updateItem := &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName: aws.String(db.TableName),
				Key: map[string]*dynamodb.AttributeValue{
					"PK": {S: aws.String(pk)},
					"SK": {S: aws.String(sk)},
				},
				UpdateExpression:          aws.String(updateExpr),
				ExpressionAttributeValues: expressionAttributeValues,
			},
		}

		transactItems = append(transactItems, updateItem)
	}

	_, err := db.Client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("failed to create update notifications: %s", err)
	}

	return nil

}
