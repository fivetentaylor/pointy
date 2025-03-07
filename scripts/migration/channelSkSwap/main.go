package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

const docPrefix = "doc#"
const ChannelPrefix = "chan#"

type OldChannel struct {
	ChannelID string `json:"channelID"`
	DocID     string `json:"docID"`
	UpdatedAt int64  `json:"updatedAt"`

	UserIDs []string
}

type dOldChannel struct {
	// doc#{docID}
	DocID string `dynamodbav:"PK"`
	// chan#{updatedAt}
	UpdatedAt string `dynamodbav:"SK"`
	// chan#{channelID}
	ChannelID string `dynamodbav:"SK1"`

	UserIDs []string `dynamodbav:"userIDs"`
}

func (item *OldChannel) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	tmp := dOldChannel{}

	if err := dynamodbattribute.UnmarshalMap(av.M, &tmp); err != nil {
		return err
	}

	var err error
	item.DocID = strings.TrimPrefix(tmp.DocID, docPrefix)
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

	return err
}

func main() {
	// Scan dynabmo DB for all objects
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := dynamo.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %+v", err)
	}

	log.SetLevel(log.DebugLevel)

	input := &dynamodb.ScanInput{
		TableName: &db.TableName,
	}

	for {
		output, err := db.Client.Scan(input)
		if err != nil {
			log.Fatalf("failed to scan table: %+v", err)
		}

		log.Infof("Scanned %d items", len(output.Items)) // output.Items

		for _, item := range output.Items {
			if item["SK"] != nil && item["SK1"] != nil {
				skA := item["SK"]
				sk1A := item["SK1"]
				if skA.S != nil && sk1A.S != nil {
					sk := *skA.S
					sk1 := *sk1A.S
					if strings.HasPrefix(sk, "chan#") && strings.HasPrefix(sk1, "chan#") {
						fmt.Printf("Updating %s to %s\n", sk, sk1)

						oldChannel := OldChannel{}
						err := dynamodbattribute.UnmarshalMap(item, &oldChannel)
						if err != nil {
							log.Fatalf("failed to unmarshal item: %+v", err)
						}

						fmt.Printf("oldChannel: %#v\n", oldChannel)

						newChannel := dynamo.Channel{
							ChannelID: oldChannel.ChannelID,
							DocID:     oldChannel.DocID,
							UpdatedAt: oldChannel.UpdatedAt,
							Type:      dynamo.ChannelTypeReviso,
							UserIDs:   oldChannel.UserIDs,
						}

						err = db.CreateChannel(&newChannel)
						if err != nil {
							log.Fatalf("failed to create channel: %+v", err)
						}

						fmt.Printf("newChannel: %#v\n", newChannel)

						deleteItem(db.Client, db.TableName, item)
					}

				}
			}
		}

		if output.LastEvaluatedKey == nil {
			log.Info("\nNo more items to scan")
			break
		}

		input.ExclusiveStartKey = output.LastEvaluatedKey
	}
}

func deleteItem(svc *dynamodb.DynamoDB, tableName string, item map[string]*dynamodb.AttributeValue) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": item["PK"], // Assuming 'PK' is the name of your partition key
			"SK": item["SK"], // Assuming 'SK' is the name of your sort key
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("failed to delete item: %+v", err)
	}
	fmt.Printf("Deleted old item: %v\n", input.Key)
}
