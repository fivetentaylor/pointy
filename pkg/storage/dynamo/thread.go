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
)

var AiThreadPrefix = "aiThrd#"

type Thread struct {
	ThreadID string `json:"threadID"`
	UserID   string `json:"userID"`

	// Attributes
	DocID     string `json:"docID"`
	Title     string `json:"title"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (item *Thread) GetKey() map[string]*dynamodb.AttributeValue {
	key := map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, item.DocID))},
		"SK": {S: aws.String(fmt.Sprintf("%s%s#%s", AiThreadPrefix, item.ThreadID, item.UserID))},
	}

	return key
}

type dThread struct {
	// PK doc#{uuid}
	PK string `dynamodbav:"PK"`
	// SK aiThrd#{threadID}#{userID}
	SK string `dynamodbav:"SK"`
	// SK1 aiThrd#{userID}#{updatedAt}
	SK1 string `dynamodbav:"SK1"`

	// Attributes
	Title string `dynamodbav:"title"`
}

func (item Thread) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	marshalItem := dThread{
		PK:    fmt.Sprintf("%s%s", DocPrefix, item.DocID),
		SK:    fmt.Sprintf("%s%s#%s", AiThreadPrefix, item.ThreadID, item.UserID),
		SK1:   fmt.Sprintf("%s%s#%d", AiThreadPrefix, item.UserID, item.UpdatedAt),
		Title: item.Title,
	}
	m, err := dynamodbattribute.MarshalMap(marshalItem)
	if err != nil {
		return err
	}
	av.M = m
	return nil
}

func (item *Thread) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dThread{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return err
	}

	log.Debugf("unmarshalling thread: %#v", tmp)

	item.DocID = strings.TrimPrefix(tmp.PK, DocPrefix)

	skParts := strings.Split(
		strings.TrimPrefix(tmp.SK, AiThreadPrefix),
		"#",
	)
	if len(skParts) != 2 {
		return fmt.Errorf("invalid sk: %s", tmp.SK)
	}
	item.ThreadID = skParts[0]

	sk1Parts := strings.Split(
		strings.TrimPrefix(tmp.SK1, AiThreadPrefix),
		"#",
	)
	if len(sk1Parts) != 2 {
		return fmt.Errorf("invalid sk1: %s", tmp.SK1)
	}

	var err error
	item.UserID = sk1Parts[0]
	item.UpdatedAt, err = strconv.ParseInt(sk1Parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid sk1: %s: %s", tmp.SK1, err)
	}

	item.Title = tmp.Title

	return nil
}

// GetThreadWithoutUserId should only be used for testing or admin
func (db *DB) GetThreadWithoutUserId(docID, threadID string) (*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND (begins_with(SK, :prefix))"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":prefix": {S: aws.String(fmt.Sprintf("%s%s", AiThreadPrefix, threadID))},
		},
		Limit: aws.Int64(1),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query thread: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, fmt.Errorf("thread not found")
	}

	var thread Thread
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &thread)
	if err != nil {
		return nil, err
	}

	if len(results.Items) > 1 {
		err = fmt.Errorf("multiple threads found")
	}

	return &thread, err
}

// Get thread
func (db *DB) GetThread(docID, threadID string) (*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND SK = :SK"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":SK": {S: aws.String(fmt.Sprintf("%s%s#", AiThreadPrefix, threadID))},
		},
		Limit: aws.Int64(1),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query thread: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, nil
	}

	var thread Thread
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &thread)
	if err != nil {
		return nil, err
	}

	if len(results.Items) > 1 {
		err = fmt.Errorf("multiple threads found")
	}

	return &thread, err
}

func (db *DB) GetThreadForUser(docID, threadID, userID string) (*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("PK = :PK AND SK = :SK"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":SK": {S: aws.String(fmt.Sprintf("%s%s#%s", AiThreadPrefix, threadID, userID))},
		},
		Limit: aws.Int64(1),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query thread: %s", err)
		return nil, err
	}

	if len(results.Items) == 0 {
		return nil, nil
	}

	var thread Thread
	err = dynamodbattribute.UnmarshalMap(results.Items[0], &thread)
	if err != nil {
		return nil, err
	}

	if len(results.Items) > 1 {
		err = fmt.Errorf("multiple threads found")
	}

	return &thread, err
}

func (db *DB) UpdateThread(updatedThread *Thread) error {
	if updatedThread.ThreadID == "" {
		return fmt.Errorf("threadID cannot be empty")
	}
	if updatedThread.DocID == "" {
		return fmt.Errorf("docID cannot be empty")
	}
	if updatedThread.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	updatedThread.UpdatedAt = time.Now().UnixNano()

	av, err := dynamodbattribute.MarshalMap(updatedThread)
	if err != nil {
		return err
	}

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":title": av["title"],
		":SK1":   av["SK1"],
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(db.TableName),
		Key:                       updatedThread.GetKey(),
		ExpressionAttributeValues: exprAttrValues,
		UpdateExpression:          aws.String("SET #SK1 = :SK1, #title = :title"),
		ExpressionAttributeNames: map[string]*string{
			"#SK1":   aws.String("SK1"),
			"#title": aws.String("title"),
		},
	}

	conditionExpression := "attribute_exists(SK)"
	updateInput.ConditionExpression = aws.String(conditionExpression)

	_, err = db.Client.UpdateItem(updateInput)
	if err != nil {
		return err
	}

	return nil
}

// Create thread
func (db *DB) CreateThread(thread *Thread) error {
	if thread.ThreadID == "" {
		thread.ThreadID = uuid.New().String()
	}
	if thread.DocID == "" {
		return fmt.Errorf("docID cannot be empty")
	}
	if thread.UserID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	thread.UpdatedAt = time.Now().UnixNano()
	log.Debugf("creating thread: %#v", thread)
	av, err := dynamodbattribute.Marshal(thread)
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

// Get all threads should only be used for testing / admin
func (db *DB) GetAllThreadsForDoc(docID string) ([]*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK1, :prefix))"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":prefix": {S: aws.String(AiThreadPrefix)},
		},
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query threads: %s", err)
		return nil, err
	}

	threads := make([]*Thread, len(results.Items))
	for i, item := range results.Items {
		thread := Thread{}
		err := dynamodbattribute.UnmarshalMap(item, &thread)
		if err != nil {
			log.Errorf("failed to unmarshal thread: %s", err)
			return nil, err
		}
		threads[i] = &thread
	}

	return threads, nil
}

func (db *DB) GetThreadsForDoc(docID string) ([]*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK1, :prefix))"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":prefix": {S: aws.String(fmt.Sprintf("%s", AiThreadPrefix))},
		},
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query threads: %s", err)
		return nil, err
	}

	threads := make([]*Thread, len(results.Items))
	for i, item := range results.Items {
		thread := Thread{}
		err := dynamodbattribute.UnmarshalMap(item, &thread)
		if err != nil {
			log.Errorf("failed to unmarshal thread: %s", err)
			return nil, err
		}
		threads[i] = &thread
	}

	return threads, nil
}

func (db *DB) GetThreadsForDocForUser(docID, userID string) ([]*Thread, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("(PK = :PK) AND (begins_with(SK1, :prefix))"),
		IndexName:              aws.String("SK1Index"),
		Select:                 aws.String("ALL_ATTRIBUTES"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK":     {S: aws.String(fmt.Sprintf("%s%s", DocPrefix, docID))},
			":prefix": {S: aws.String(fmt.Sprintf("%s%s#", AiThreadPrefix, userID))},
		},
		ScanIndexForward: aws.Bool(false),
	}

	results, err := db.Client.Query(input)
	if err != nil {
		log.Errorf("failed to query threads: %s", err)
		return nil, err
	}

	threads := make([]*Thread, len(results.Items))
	for i, item := range results.Items {
		thread := Thread{}
		err := dynamodbattribute.UnmarshalMap(item, &thread)
		if err != nil {
			log.Errorf("failed to unmarshal thread: %s", err)
			return nil, err
		}
		threads[i] = &thread
	}

	return threads, nil
}

func (db *DB) DeleteThreadsForDoc(docID, userID string) ([]*Thread, error) {
	threads, err := db.GetThreadsForDocForUser(docID, userID)
	if err != nil {
		return nil, err
	}

	for _, thread := range threads {
		err = db.DeleteThread(thread)
		if err != nil {
			return nil, err
		}
	}

	return threads, nil
}

func (db *DB) DeleteThread(thread *Thread) error {
	_, err := db.DeleteMessagesForThread(thread.ThreadID)
	if err != nil {
		return err
	}

	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key:       thread.GetKey(),
	}

	_, err = db.Client.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}
