package dynamo

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/teamreviso/code/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
)

const UserPreferencePrefix = "userPref#"
const DocNotificationPreferencePrefix = "docNotifPref#"

var DefaultDocPreference = &models.DocumentPreference{
	EnableFirstOpenNotifications:  true,
	EnableAllCommentNotifications: false,
	EnableMentionNotifications:    true,
	EnableDmNotifications:         true,
}

var ValidUnreadActivityFrequencyMinutes = []int32{2, 5, 10}

var DefaultUserPreference = &models.UserPreference{
	EnableActivityNotifications:    true,
	UnreadActivityFrequencyMinutes: 2,
}

type UserPreference struct {
	// PK and SK
	UserID string `json:"userID"`

	// Attributes
	Preference *models.UserPreference `json:"preference"`
}

type dUserNotificationPreference struct {
	// PK: userPref#userID
	PK string `dynamodbav:"PK" json:"PK"`
	// SK: userID
	SK string `dynamodbav:"SK" json:"SK"`

	Preference string `dynamodbav:"preference" json:"preference"`
}

func (n *UserPreference) Key() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", UserPreferencePrefix, n.UserID))},
		"SK": {S: aws.String(n.UserID)},
	}
}

func (n UserPreference) toDynamo() (dUserNotificationPreference, error) {
	dn := dUserNotificationPreference{
		PK: fmt.Sprintf("%s%s", DocNotificationPreferencePrefix, n.UserID),
		SK: n.UserID,
	}

	bts, err := protojson.Marshal(n.Preference)
	if err != nil {
		return dn, fmt.Errorf("failed to marshal payload: %s", err)
	}

	dn.Preference = string(bts)
	return dn, nil
}

func (n *UserPreference) fromDynamo(dn dUserNotificationPreference) error {
	n.UserID = dn.SK

	p := &models.UserPreference{}
	if err := protojson.Unmarshal([]byte(dn.Preference), p); err != nil {
		return fmt.Errorf("failed to unmarshal preference: %s", err)
	}
	n.Preference = p

	return nil
}

func (n UserPreference) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
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

func (n *UserPreference) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dUserNotificationPreference{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := n.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (db *DB) GetUserPreference(userID string) (*UserPreference, error) {
	userPref := UserPreference{
		UserID: userID,
	}

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       userPref.Key(),
	}

	result, err := db.Client.GetItem(getInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %s", err)
	}

	if result.Item == nil {
		userPref.Preference = DefaultUserPreference
		return &userPref, nil
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &userPref); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %s", err)
	}

	return &userPref, nil
}

func (db *DB) UpsertUserPreference(n *UserPreference) error {
	if n.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	if n.Preference == nil {
		return fmt.Errorf("preference cannot be empty")
	}

	var valid bool
	for _, v := range ValidUnreadActivityFrequencyMinutes {
		if v == n.Preference.UnreadActivityFrequencyMinutes {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid value for unreadActivityFrequencyMinutes")
	}

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
		return fmt.Errorf("UpsertUserPreferences(%+v) (%+v): %w", n, input, err)
	}

	return nil
}

type DocPreference struct {
	// PK
	UserID string `json:"userID"`

	// SK
	DocID string `json:"docID"`

	// Attributes
	Preference *models.DocumentPreference `json:"preference"`
}

func (n *DocPreference) Key() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", DocNotificationPreferencePrefix, n.UserID))},
		"SK": {S: aws.String(n.DocID)},
	}
}

func (n DocPreference) toDynamo() (dDocNotificationPreference, error) {
	dn := dDocNotificationPreference{
		PK: fmt.Sprintf("%s%s", DocNotificationPreferencePrefix, n.UserID),
		SK: n.DocID,
	}

	bts, err := protojson.Marshal(n.Preference)
	if err != nil {
		return dn, fmt.Errorf("failed to marshal payload: %s", err)
	}

	dn.Preference = string(bts)
	return dn, nil
}

func (n *DocPreference) fromDynamo(dn dDocNotificationPreference) error {
	n.UserID = strings.TrimPrefix(dn.PK, DocNotificationPreferencePrefix)
	n.DocID = dn.SK

	p := &models.DocumentPreference{}
	if err := protojson.Unmarshal([]byte(dn.Preference), p); err != nil {
		return fmt.Errorf("failed to unmarshal preference: %s", err)
	}
	n.Preference = p

	return nil
}

type dDocNotificationPreference struct {
	// PK: docNotifPref#userID
	PK string `dynamodbav:"PK" json:"PK"`
	// SK: docID
	SK string `dynamodbav:"SK" json:"SK"`

	Preference string `dynamodbav:"preference" json:"preference"`
}

func (n DocPreference) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
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

func (n *DocPreference) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dDocNotificationPreference{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := n.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (db *DB) GetDocNotificationPreference(userID, docID string) (*DocPreference, error) {
	docNotif := DocPreference{
		UserID: userID,
		DocID:  docID,
	}

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       docNotif.Key(),
	}

	result, err := db.Client.GetItem(getInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %s", err)
	}

	if result.Item == nil {
		docNotif.Preference = DefaultDocPreference
		return &docNotif, nil
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &docNotif); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %s", err)
	}

	return &docNotif, nil
}

func (db *DB) UpsertDocNotificationPreference(n *DocPreference) error {
	if n.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	if n.DocID == "" {
		return fmt.Errorf("docID cannot be empty")
	}

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
		return fmt.Errorf("CreateOrUpdateDocNotificationPreference(%+v) (%+v): %w", n, input, err)
	}

	return nil
}
