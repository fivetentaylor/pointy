package dynamo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
	"github.com/fivetentaylor/pointy/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

const MsgPrefix = "msg#"
const MainMsgChain = "main"

var protoJsonUnmarshaler = protojson.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

type MessageLifecycleStage int

const (
	MessageLifecycleStageUnknown MessageLifecycleStage = iota
	MessageLifecycleStagePending
	MessageLifecycleStageRevised
	MessageLifecycleStageCompleted
	MessageLifecycleStageRevising
)

func (m MessageLifecycleStage) String() string {
	switch m {
	case MessageLifecycleStagePending:
		return "pending"
	case MessageLifecycleStageRevised:
		return "revised"
	case MessageLifecycleStageCompleted:
		return "completed"
	case MessageLifecycleStageRevising:
		return "revising"
	default:
		return "unknown"
	}
}

type MessageType string

const (
	MessageTypeMention MessageType = "mention"
	MessageTypeDM      MessageType = "dm"
	MessageTypeComment MessageType = "comment"
)

type Message struct {
	// PK
	ContainerID string `json:"containerID"`

	// SK
	MessageID string `json:"messageID"`
	// SK1
	Chain     string `json:"chain"`
	CreatedAt int64  `json:"createdAt"`

	// Attributes
	UserID          string                  `json:"userID"`
	AuthorID        string                  `json:"authorID"`
	ChannelID       string                  `json:"channelID"`
	Content         string                  `json:"content"`
	LifecycleStage  MessageLifecycleStage   `json:"lifecycleStage"`
	LifecycleReason string                  `json:"lifecycleReason"`
	Attachments     *models.AttachmentList  `json:"attachments"`
	AIContent       *models.AIContent       `json:"aiContent"`
	MessageMetadata *models.MessageMetadata `json:"messageMetadata"`
	Hidden          bool                    `json:"hidden"`

	// Reply
	ReplyCount int `json:"replyCount"`

	// Relationships
	DocID               string   `json:"docID"`
	ParentContainerID   *string  `json:"parentContainerID"`
	ForkedFromMessageID *string  `json:"forkedFrom"`
	ForkedMessages      []string `json:"forkedMessages"`
	ReplyingUserIds     []string `json:"replyingUserIds"`
	MentionedUserIds    []string `json:"mentionedUserIds"`
}

func (m *Message) GetContainerID() string {
	return fmt.Sprintf("%s%s", MsgPrefix, m.MessageID)
}

func (m *Message) CreatedAtTimestamp() time.Time {
	return time.UnixMicro(m.CreatedAt / 1000)
}

func (m *Message) GetMessageType(channel *Channel, userID string) models.CommentType {
	if channel.Type == ChannelTypeDirect {
		return models.CommentType_DirectMessage
	} else if utils.Contains(m.MentionedUserIds, userID) {
		return models.CommentType_Mention
	} else if m.HasSelection() {
		return models.CommentType_Comment
	} else {
		return models.CommentType_Message
	}
}

func (m *Message) FullContent() string {
	var b strings.Builder

	for _, attachment := range m.Attachments.Attachments {
		switch v := attachment.Value.(type) {
		case *models.Attachment_Content:
			b.WriteString(v.Content.Text)
			b.WriteString("\n")
		}
	}

	b.WriteString(m.Content)
	b.WriteString("\n")

	return b.String()
}

func (m *Message) Reload(db *DB) (*Message, error) {
	return db.GetMessage(m.ContainerID, m.MessageID)
}

func (m *Message) HasSelection() bool {
	if m.Attachments == nil || len(m.Attachments.Attachments) == 0 {
		return false
	}

	for _, attachment := range m.Attachments.Attachments {
		switch attachment.Value.(type) {
		case *models.Attachment_Document:
			return true
		}
	}

	return false
}

// GetContentSize is used to see if the message has change after a stream delta merge
func (m *Message) GetContentSize() int {
	contentSize := len(m.Content)
	for _, attachment := range m.Attachments.Attachments {
		switch v := attachment.Value.(type) {
		case *models.Attachment_Document:
			contentSize += len(v.Document.Content)
		case *models.Attachment_Revision:
			contentSize += len(v.Revision.Updated)
		case *models.Attachment_Suggestion:
			contentSize += len(v.Suggestion.Content)
		case *models.Attachment_Content:
			contentSize += len(v.Content.Text)
		}
	}
	if m.AIContent != nil {
		contentSize += len(m.AIContent.ConcludingMessage)
		contentSize += len(m.AIContent.Feedback)
	}

	return contentSize
}

func (m *Message) toDynamo() (dMessage, error) {
	marshalItem := dMessage{
		PK:                  fmt.Sprintf(m.ContainerID),
		SK:                  fmt.Sprintf("%s%s", MsgPrefix, m.MessageID),
		SK1:                 fmt.Sprintf("%s%s@%d", MsgPrefix, m.Chain, m.CreatedAt),
		Content:             m.Content,
		DocID:               m.DocID,
		UserID:              m.UserID,
		AuthorID:            m.AuthorID,
		ChannelID:           m.ChannelID,
		LifecycleStage:      int(m.LifecycleStage),
		LifecycleReason:     m.LifecycleReason,
		ReplyCount:          m.ReplyCount,
		ParentContainerID:   m.ParentContainerID,
		ForkedFromMessageID: m.ForkedFromMessageID,
		ForkedMessages:      m.ForkedMessages,
		ReplyingUserIds:     m.ReplyingUserIds,
		MentionedUserIds:    m.MentionedUserIds,
		Hidden:              m.Hidden,
	}

	abts, err := protojson.Marshal(m.Attachments)
	if err != nil {
		return marshalItem, fmt.Errorf("failed to marshal attachments: %s", err)
	}

	marshalItem.Attachments = string(abts)

	aibts, err := protojson.Marshal(m.AIContent)
	if err != nil {
		return marshalItem, fmt.Errorf("failed to marshal aiContent: %s", err)
	}

	marshalItem.AIContent = string(aibts)

	mobts, err := protojson.Marshal(m.MessageMetadata)
	if err != nil {
		return marshalItem, fmt.Errorf("failed to marshal messageMetadata: %s", err)
	}

	marshalItem.MessageMetadata = string(mobts)

	return marshalItem, nil
}

func (m *Message) fromDynamo(tmp dMessage) error {
	var err error
	m.DocID = tmp.DocID
	m.ContainerID = tmp.PK
	m.MessageID = strings.TrimPrefix(tmp.SK, MsgPrefix)
	m.Content = tmp.Content
	m.UserID = tmp.UserID
	m.AuthorID = tmp.AuthorID
	m.ChannelID = tmp.ChannelID
	m.LifecycleStage = MessageLifecycleStage(tmp.LifecycleStage)
	m.LifecycleReason = tmp.LifecycleReason
	m.ReplyCount = tmp.ReplyCount
	m.ParentContainerID = tmp.ParentContainerID
	m.ForkedFromMessageID = tmp.ForkedFromMessageID
	m.ForkedMessages = tmp.ForkedMessages
	m.ReplyingUserIds = tmp.ReplyingUserIds
	if m.ReplyingUserIds == nil {
		m.ReplyingUserIds = []string{}
	}
	m.MentionedUserIds = tmp.MentionedUserIds
	if m.MentionedUserIds == nil {
		m.MentionedUserIds = []string{}
	}

	sk1Parts := strings.Split(
		strings.TrimPrefix(tmp.SK1, MsgPrefix),
		"@",
	)

	if len(sk1Parts) != 2 {
		return fmt.Errorf("invalid sk1: %s", tmp.SK1)
	}

	m.Chain = sk1Parts[0]
	m.CreatedAt, err = strconv.ParseInt(sk1Parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse createdAt: %s", err)
	}

	abts := []byte(tmp.Attachments)
	m.Attachments = &models.AttachmentList{}
	err = protoJsonUnmarshaler.Unmarshal(abts, m.Attachments)
	if err != nil {
		return fmt.Errorf("failed to unmarshal attachments: %s\n%s\n%q", err, m.GetKey(), tmp.Attachments)
	}

	if tmp.AIContent != "" {
		aibts := []byte(tmp.AIContent)
		if m.AIContent == nil {
			m.AIContent = &models.AIContent{}
		}
		err = protoJsonUnmarshaler.Unmarshal(aibts, m.AIContent)
		if err != nil {
			return fmt.Errorf("failed to unmarshal aiContent: %s", err)
		}
	}

	m.MessageMetadata = &models.MessageMetadata{}
	if tmp.MessageMetadata != "" {
		mobts := []byte(tmp.MessageMetadata)
		if m.MessageMetadata == nil {
			m.MessageMetadata = &models.MessageMetadata{}
		}
		err = protoJsonUnmarshaler.Unmarshal(mobts, m.MessageMetadata)
		if err != nil {
			return fmt.Errorf("failed to unmarshal messageMetadata: %s", err)
		}
	}
	m.Hidden = tmp.Hidden

	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	marshalItem, err := m.toDynamo()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal internal message: %s", err)
	}

	bts, err := json.Marshal(marshalItem)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal internal message: %s", err)
	}
	return bts, nil
}

func (m *Message) UnmarshalJSON(b []byte) error {
	tmp := dMessage{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := m.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (m *Message) GetKey() map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(m.ContainerID)},
		"SK": {S: aws.String(fmt.Sprintf("%s%s", MsgPrefix, m.MessageID))},
	}

	return key
}

type dMessage struct {
	// chan#{uuid} or msg#{uuid}
	PK string `dynamodbav:"PK" json:"PK"`
	// msg#{uuid}
	SK string `dynamodbav:"SK" json:"SK"`
	// msg#main@{createdAt} or msg#msg#{uuid}@{createdAt}
	SK1 string `dynamodbav:"SK1" json:"SK1"`

	UserID          string `dynamodbav:"userID" json:"userID"`
	AuthorID        string `dynamodbav:"authorID" json:"authorID"`
	ChannelID       string `dynamodbav:"channelID" json:"channelID"`
	Content         string `dynamodbav:"content" json:"content"`
	LifecycleStage  int    `dynamodbav:"lifecycleStage" json:"lifecycleStage"`
	LifecycleReason string `dynamodbav:"lifecycleReason" json:"lifecycleReason"`
	Attachments     string `dynamodbav:"attachments" json:"attachments"`
	AIContent       string `dynamodbav:"aiContent" json:"aiContent"`
	MessageMetadata string `dynamodbav:"messageMetadata" json:"messageMetadata"`
	Hidden          bool   `dynamodbav:"hidden" json:"hidden"`

	ReplyCount int `dynamodbav:"replyCount" json:"replyCount"`

	DocID               string   `dynamodbav:"docID" json:"docID"`
	ParentContainerID   *string  `dynamodbav:"parentContainerID" json:"parentContainerID"`
	ForkedFromMessageID *string  `dynamodbav:"forkedFrom" json:"forkedFrom"`
	ForkedMessages      []string `dynamodbav:"forkedMessages" json:"forkedMessages"`
	ReplyingUserIds     []string `dynamodbav:"replyingUserIds" json:"replyingUserIds"`
	MentionedUserIds    []string `dynamodbav:"mentionedUserIds" json:"mentionedUserIds"`
}

func (m Message) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	marshalItem, err := m.toDynamo()
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

func (m *Message) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dMessage{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal internal message: %s", err)
	}

	if err := m.fromDynamo(tmp); err != nil {
		return fmt.Errorf("failed to unmarshal message: %s", err)
	}

	return nil
}

func (db *DB) CreateMessage(message *Message) error {
	if message.DocID == "" {
		return fmt.Errorf("docID is required")
	}
	if message.ContainerID == "" {
		return fmt.Errorf("containerID is required")
	}
	if message.AuthorID == "" {
		return fmt.Errorf("authorID is required")
	}
	if message.UserID == "" {
		return fmt.Errorf("userID is required")
	}
	if message.ChannelID == "" {
		return fmt.Errorf("channelID is required")
	}
	if message.Chain == "" {
		message.Chain = MainMsgChain
	}
	if message.MessageID == "" {
		message.MessageID = uuid.NewString()
	}
	if message.ReplyingUserIds == nil {
		message.ReplyingUserIds = []string{}
	}

	if strings.HasPrefix(MsgPrefix, message.ContainerID) && message.ParentContainerID == nil {
		return fmt.Errorf("parentContainerID is required for threadMessages")
	}

	message.CreatedAt = time.Now().UnixNano()

	if message.Attachments == nil {
		message.Attachments = &models.AttachmentList{}
	}

	if message.AIContent == nil {
		message.AIContent = &models.AIContent{}
	}

	if message.MessageMetadata == nil {
		message.MessageMetadata = &models.MessageMetadata{}
	}

	av, err := dynamodbattribute.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	log.Info("[dydb] creating message", "message", message)

	transactItems := []*dynamodb.TransactWriteItem{}

	conditionExpression := "attribute_not_exists(SK)"
	putItem := &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:           aws.String(db.TableName),
			ConditionExpression: aws.String(conditionExpression),
			Item:                av.M,
		},
	}
	transactItems = append(transactItems, putItem)

	_, err = db.Client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("failed to create message and update count: %s", err)
	}

	return nil
}

func (db *DB) UpdateMessage(msg *Message) error {
	// Check for required fields
	if msg.MessageID == "" {
		return fmt.Errorf("messageID is required")
	}
	if msg.CreatedAt == 0 {
		return fmt.Errorf("createdAt is required")
	}

	// Marshal the Message to a map[string]*dynamodb.AttributeValue
	av, err := dynamodbattribute.MarshalMap(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	// Define the update expression, attribute names and values
	updateExpression := "SET #replyCount = :replyCount, #content = :content,  #lifecycleStage = :lifecycleStage, #attachments = :attachments, #aiContent = :aiContent, #messageMetadata = :messageMetadata, #SK1 = :sk1, #forkedMessages = :forkedMessages"
	exprAttrNames := map[string]*string{
		"#replyCount":      aws.String("replyCount"),
		"#content":         aws.String("content"),
		"#lifecycleStage":  aws.String("lifecycleStage"),
		"#attachments":     aws.String("attachments"),
		"#aiContent":       aws.String("aiContent"),
		"#messageMetadata": aws.String("messageMetadata"),
		"#SK1":             aws.String("SK1"),
		"#forkedMessages":  aws.String("forkedMessages"),
	}
	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":replyCount":      av["replyCount"],
		":content":         av["content"],
		":lifecycleStage":  av["lifecycleStage"],
		":attachments":     av["attachments"],
		":aiContent":       av["aiContent"],
		":messageMetadata": av["messageMetadata"],
		":sk1":             av["SK1"],
		":forkedMessages":  av["forkedMessages"],
	}

	// Create the update item input
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.TableName),
		Key:                       msg.GetKey(),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Perform the update operation
	_, err = db.Client.UpdateItem(updateItemInput)
	if err != nil {
		return fmt.Errorf("failed to update item: %s", err)
	}

	return nil
}

func (db *DB) GetChannelMessages(channelID string) ([]*Message, error) {
	log.Infof("[dydb] GetChannelMessages")
	return db.getMessages(
		fmt.Sprintf("%s%s", ChannelPrefix, channelID),
		"main",
	)
}

func (db *DB) GetThreadMessages(containerID, msgID string) ([]*Message, error) {
	log.Info("[dydb] GetThreadMessages", "containerID", containerID, "msgID", msgID)
	initialMsg, err := db.GetMessage(containerID, msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial message: %s", err)
	}

	msgs, err := db.getMessages(
		fmt.Sprintf("%s%s", MsgPrefix, msgID),
		"main",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %s", err)
	}

	return append([]*Message{initialMsg}, msgs...), nil
}

func (db *DB) GetRepliesToMessage(msgID string) ([]*Message, error) {
	log.Infof("[dydb] GetRepliesToMessage")

	msgs, err := db.getMessages(
		fmt.Sprintf("%s%s", MsgPrefix, msgID),
		"main",
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (db *DB) GetAlternateRepliesToMessage(msgID, chain string) ([]*Message, error) {
	log.Infof("[dydb] GetAlternateRepliesToMessage")

	msgs, err := db.getMessages(
		fmt.Sprintf("%s%s", MsgPrefix, msgID),
		chain,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (db *DB) GetMessagesForThread(threadID string) ([]*Message, error) {
	log.Infof("[dydb] GetMessagesForThread")
	return db.getMessages(
		fmt.Sprintf("%s%s", AiThreadPrefix, threadID),
		"main",
	)
}

func (db *DB) DeleteMessagesForThread(threadID string) ([]*Message, error) {
	messages, err := db.GetMessagesForThread(threadID)
	if err != nil {
		return nil, err
	}

	for _, msg := range messages {
		err := db.DeleteMessage(msg)
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}

func (db *DB) DeleteMessage(msg *Message) error {
	log.Infof("[dydb] DeleteMessage")
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       msg.GetKey(),
	}

	_, err := db.Client.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getMessages(containerID, chain string) ([]*Message, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK1, :prefix))"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(containerID)},
			":prefix": {S: aws.String(fmt.Sprintf("%s%s", MsgPrefix, chain))},
		},
		ScanIndexForward: aws.Bool(true),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		return nil, stackerr.Errorf("failed to query messages: %s", err)
	}

	messages := make([]*Message, len(results.Items))
	for i, item := range results.Items {
		msg := Message{}
		err := dynamodbattribute.UnmarshalMap(item, &msg)
		if err != nil {
			log.Errorf("failed to unmarshal message map: %s", err)
			return nil, err
		}
		messages[i] = &msg
	}

	return messages, nil
}

func (db *DB) GetAiThreadMessage(threadID, messageID string) (*Message, error) {
	return db.GetMessage(
		fmt.Sprintf("%s%s", AiThreadPrefix, threadID),
		messageID,
	)
}

func (db *DB) GetMessage(containerID, messageID string) (*Message, error) {
	log.Info("[dydb] GetMessage", "containerID", containerID, "messageID", messageID)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("(PK = :PK) AND SK = :SK"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(containerID)},
			":SK": {S: aws.String(fmt.Sprintf("%s%s", MsgPrefix, messageID))},
		},
		Limit: aws.Int64(1),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query message %+v: %s", input, err)
		return nil, stackerr.Errorf("failed to query message: %s", err)
	}

	if len(results.Items) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	var msg Message
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %s", err)
	}

	if len(results.Items) > 1 {
		err = fmt.Errorf("multiple messages found")
	}

	return &msg, err
}

// GetRandomMessage returns a random message
// ONLY SHOULD BE USED IN DEVELOPMENT
func (db *DB) GetRandomMessage() (*Message, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(db.TableName),
		FilterExpression: aws.String("begins_with(PK, :PKPrefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PKPrefix": {S: aws.String(MsgPrefix)},
		},
		Limit: aws.Int64(int64(rand.Intn(200) + 1)),
	}

	results, err := db.Client.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for messages: %s", err)
	}

	// Check if we have any items
	numItems := len(results.Items)
	if numItems == 0 {
		return nil, fmt.Errorf("no messages found")
	}

	// Randomly select one item from the results
	randomIndex := rand.Intn(numItems)
	var msg Message
	err = dynamodbattribute.UnmarshalMap(results.Items[randomIndex], &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message %#v:  %s", results.Items[randomIndex], err)
	}

	return &msg, nil
}
