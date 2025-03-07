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
	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

var ChannelPrefix = "chan#"

type ChannelType int

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeReviso
	ChannelTypeDirect
	ChannelTypeGeneral
)

func (item ChannelType) String() string {
	switch item {
	case ChannelTypeReviso:
		return "reviso"
	case ChannelTypeDirect:
		return "direct"
	case ChannelTypeGeneral:
		return "general"
	default:
		return "unknown"
	}
}

type Channel struct {
	ChannelID string      `json:"channelID"`
	DocID     string      `json:"docID"`
	UpdatedAt int64       `json:"updatedAt"`
	Type      ChannelType `json:"type"`

	UserIDs []string
}

func (item Channel) GetContainerID() string {
	return fmt.Sprintf("%s%s", ChannelPrefix, item.ChannelID)
}

func (item Channel) ContainerID() string {
	return fmt.Sprintf("%s%s", ChannelPrefix, item.ChannelID)
}

func (item *Channel) GetKey() map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, item.DocID))},
		"SK": {S: aws.String(fmt.Sprintf("%s%s", ChannelPrefix, item.ChannelID))},
	}

	return key
}

type dChannel struct {
	// doc#{docID}
	DocID string `dynamodbav:"PK"`
	// chan#{channelID}
	ChannelID string `dynamodbav:"SK"`
	// chan#{updatedAt}
	UpdatedAt string   `dynamodbav:"SK1"`
	Type      int      `dynamodbav:"type"`
	UserIDs   []string `dynamodbav:"userIDs"`
}

func (item Channel) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	marshalItem := dChannel{
		DocID:     fmt.Sprintf("%s%s", DocPrefix, item.DocID),
		UpdatedAt: fmt.Sprintf("%s%d", ChannelPrefix, item.UpdatedAt),
		ChannelID: fmt.Sprintf("%s%s", ChannelPrefix, item.ChannelID),
		Type:      int(item.Type),
		UserIDs:   item.UserIDs,
	}
	m, err := dynamodbattribute.MarshalMap(marshalItem)
	if err != nil {
		return err
	}
	av.M = m
	return nil
}

func (item *Channel) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dChannel{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return err
	}

	var err error
	item.DocID = strings.TrimPrefix(tmp.DocID, DocPrefix)
	item.ChannelID = strings.TrimPrefix(tmp.ChannelID, ChannelPrefix)
	item.UserIDs = tmp.UserIDs
	item.UpdatedAt, err = strconv.ParseInt(
		strings.TrimPrefix(
			tmp.UpdatedAt,
			ChannelPrefix,
		),
		10, 64)

	if item.UserIDs == nil {
		item.UserIDs = []string{}
	}

	item.Type = ChannelType(tmp.Type)

	return err
}

// Get channel
func (db *DB) GetChannel(docID, channelID string) (*Channel, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND SK = :SK"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":SK": {S: aws.String(fmt.Sprintf("%s%s", ChannelPrefix, channelID))},
		},
		Limit: aws.Int64(1),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query channel: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, fmt.Errorf("channel not found")
	}

	var channel Channel
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &channel)
	if err != nil {
		return nil, err
	}

	if len(results.Items) > 1 {
		err = fmt.Errorf("multiple channels found")
	}

	return &channel, err
}

// Get all document channels
func (db *DB) GetDocumentChannels(docID string) ([]*Channel, error) {
	return db.getChannels(docID, ChannelPrefix)
}

// Delete all document channels
func (db *DB) DeleteDocumentChannels(docID string) error {
	channels, err := db.GetDocumentChannels(docID)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		err = db.DeleteChannel(channel)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create channel
func (db *DB) CreateChannel(channel *Channel) error {
	if channel.ChannelID == "" {
		channel.ChannelID = uuid.NewString()
	}
	if channel.UserIDs == nil {
		channel.UserIDs = []string{}
	}
	if channel.Type == ChannelTypeUnknown {
		return fmt.Errorf("unknown channel type")
	}
	channel.UpdatedAt = time.Now().UnixNano()

	log.Infof("creating channel: %#v", channel)
	av, err := dynamodbattribute.Marshal(channel)
	if err != nil {
		return err
	}

	conditionExpression := "attribute_not_exists(SK)"
	putInput := &dynamodb.PutItemInput{
		TableName:           aws.String(db.TableName),
		ConditionExpression: aws.String(conditionExpression),
		Item:                av.M,
	}

	_, err = db.Client.PutItem(putInput)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteChannel(channel *Channel) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       channel.GetKey(),
	}

	_, err := db.Client.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateChannel(updatedChannel *Channel) error {
	if updatedChannel.ChannelID == "" {
		return fmt.Errorf("channelID cannot be empty")
	}
	if updatedChannel.DocID == "" {
		return fmt.Errorf("docID cannot be empty")
	}

	updatedChannel.UpdatedAt = time.Now().UnixNano()

	av, err := dynamodbattribute.MarshalMap(updatedChannel)
	if err != nil {
		return err
	}

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":type":    av["type"],
		":userIDs": av["userIDs"],
		":SK1":     av["SK1"],
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.TableName),
		Key:                       updatedChannel.GetKey(),
		ExpressionAttributeValues: exprAttrValues,
		UpdateExpression:          aws.String("SET #SK1 = :SK1, #userIDs = :userIDs, #type = :type"),
		ExpressionAttributeNames: map[string]*string{
			"#SK1":     aws.String("SK1"),
			"#userIDs": aws.String("userIDs"),
			"#type":    aws.String("type"),
		},
	}

	_, err = db.Client.UpdateItem(updateInput)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChannelsForAllPairs(docID string, userIDs []string) (int, error) {
	currentChannels, err := db.GetDocumentChannels(docID)
	if err != nil {
		return 0, fmt.Errorf("failed to get current channels: %w", err)
	}

	var count int
	for i := 0; i < len(userIDs); i++ {
		for j := i + 1; j < len(userIDs); j++ {
			created, err := db.ensureChannel(docID, currentChannels, userIDs[i], userIDs[j])
			if err != nil {
				return count, fmt.Errorf("failed to ensure channel: %w", err)
			}
			if created {
				count++
			}
		}
	}

	return count, nil
}

func (db *DB) ensureChannel(
	docID string,
	channels []*Channel,
	user1, user2 string,
) (bool, error) {
	if user1 == user2 {
		return false, nil
	}
	for _, t := range channels {
		if utils.Contains(t.UserIDs, user1) && utils.Contains(t.UserIDs, user2) {
			return false, nil
		}
	}

	// Add new channel
	channel := &Channel{
		DocID:   docID,
		UserIDs: []string{user1, user2},
	}

	err := db.CreateChannel(channel)
	if err != nil {
		return false, fmt.Errorf("error creating channel: %w", err)
	}

	return true, nil
}

func (db *DB) getChannels(docID, prefix string) ([]*Channel, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK, :prefix))"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":prefix": {S: aws.String(prefix)},
		},
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query channels: %s", err)
		return nil, err
	}

	channels := make([]*Channel, len(results.Items))
	for i, item := range results.Items {
		channel := Channel{}
		err := dynamodbattribute.UnmarshalMap(item, &channel)
		if err != nil {
			log.Errorf("failed to unmarshal channel: %s", err)
			return nil, err
		}
		channels[i] = &channel
	}

	return channels, nil
}
