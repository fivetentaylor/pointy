package dynamo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
)

const DocPrefix = "doc#"

type DB struct {
	TableName string
	Client    *dynamodb.DynamoDB
}

func NewDB() (*DB, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 1, // Set the timeout to 10 seconds
	}

	awsConfig := aws.Config{
		HTTPClient: httpClient,
	}

	dynamodb_url, exists := os.LookupEnv("AWS_DYNAMODB_URL")
	if exists {
		awsConfig.Endpoint = aws.String(dynamodb_url)
	}

	dynamodb_table, exists := os.LookupEnv("AWS_DYNAMODB_TABLE")
	if !exists {
		return nil, fmt.Errorf("AWS_DYNAMODB_TABLE not set in environment variables")
	}

	slog.Info("Connecting to dynamodb tabl", "table", dynamodb_table)

	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	slog.Info("Connected to dynamodb url", "url", dynamodb_url)
	return &DB{
		TableName: dynamodb_table,
		Client:    dynamodb.New(sess),
	}, nil

}

func (db *DB) CreateTestTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(db.TableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK1"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK2"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK3"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK4"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK5"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GSI1SK"), // Include this only if you've added GSI1SK to the GSI key schema
				AttributeType: aws.String("S"),
			},
		},
		LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			{
				IndexName: aws.String("SK1Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("SK1"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
			{
				IndexName: aws.String("SK2Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("SK2"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
			{
				IndexName: aws.String("SK3Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("SK3"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
			{
				IndexName: aws.String("SK4Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("SK4"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
			{
				IndexName: aws.String("SK5Index"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("SK5"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("GSI1PK"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("GSI1SK"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"), // You may choose KEYS_ONLY, INCLUDE, or ALL based on your needs
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := db.Client.CreateTable(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "ResourceInUseException" {
			slog.Info("CreateTables errored, but that's okay. Table already exists 👍")
		} else {
			log.Errorf("failed to create table: %+v", err)
			return err
		}
	}

	// Wait for table to be active
	for {
		desc, err := db.Client.DescribeTable(&dynamodb.DescribeTableInput{
			TableName: aws.String(db.TableName),
		})
		if err != nil {
			log.Errorf("Describe table failed: %v", err)
			return err
		}

		if *desc.Table.TableStatus == "ACTIVE" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	slog.Info("Table created successfully!")
	return nil
}

// Helper function to check if the error is a conditional check failed error
func isConditionalCheckFailedError(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException
	}
	return false
}

func (db *DB) PutUniqueItem(item map[string]*dynamodb.AttributeValue) error {
	conditionExpression := "attribute_not_exists(SK)"

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(db.TableName),
		ConditionExpression: aws.String(conditionExpression),
		Item:                item,
	}

	_, err := db.Client.PutItem(input)
	if err != nil {
		log.Printf("failed to put item to DynamoDB: %v", err)
		return err
	}

	return nil
}

func (db *DB) GetItem(pk, sk string) (map[string]*dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
	}

	result, err := db.Client.GetItem(input)
	if err != nil {
		log.Printf("failed to get item from DynamoDB: %v", err)
		return nil, err
	}

	if result.Item == nil {
		log.Error("no item found with the specified PK and SK")
		return nil, fmt.Errorf("no item found with PK: %s and SK: %s", pk, sk)
	}

	return result.Item, nil
}

func ItemToJson(item map[string]*dynamodb.AttributeValue) ([]byte, error) {
	var obj interface{} // or replace with your specific struct if you have one
	err := dynamodbattribute.UnmarshalMap(item, &obj)
	if err != nil {
		log.Printf("failed to unmarshal item: %v", err)
		return nil, err
	}

	itemJSON, err := json.Marshal(obj)
	if err != nil {
		log.Printf("failed to marshal item to JSON: %v", err)
		return nil, err
	}

	return itemJSON, nil
}

func (db *DB) ScanTable() ([]map[string]*dynamodb.AttributeValue, error) {
	var items []map[string]*dynamodb.AttributeValue

	input := &dynamodb.ScanInput{
		TableName: &db.TableName,
	}

	for {
		result, err := db.Client.Scan(input)
		if err != nil {
			return nil, err
		}

		items = append(items, result.Items...)

		if result.LastEvaluatedKey == nil {
			break
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return items, nil
}
