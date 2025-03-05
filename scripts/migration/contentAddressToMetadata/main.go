package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func main() {
	// Scan dynabmo DB for all objects
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	log.SetLevel(log.DebugLevel)

	db, err := dynamo.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %+v", err)
	}

	input := &dynamodb.ScanInput{
		TableName: &db.TableName,
	}

	for {
		output, err := db.Client.Scan(input)
		if err != nil {
			log.Fatalf("failed to scan table: %+v", err)
		}

		log.Info("Scanned %d items", len(output.Items)) // output.Items

		for _, item := range output.Items {
			fmt.Printf("\n")
			sk, ok := item["SK"]
			if !ok {
				continue
			}
			skStr := *sk.S
			fmt.Printf("%s", *item["SK"].S)
			if !strings.HasPrefix(skStr, "msg#") {
				continue
			}
			fmt.Printf(" | %s", *item["PK"].S)

			contentAddress, ok := item["contentAddress"]
			if !ok {
				fmt.Printf(" | %s", "no content address")
				continue
			}

			if contentAddress.S == nil {
				fmt.Printf(" | %s", "no content address")
				continue
			}

			fmt.Printf(" | %s", *contentAddress.S)

			msg := dynamo.Message{}
			err = dynamodbattribute.UnmarshalMap(item, &msg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//set contentAddress on Metadata
			if msg.MessageMetadata == nil {
				msg.MessageMetadata = &models.MessageMetadata{
					ContentAddress: *contentAddress.S,
				}
			} else {
				msg.MessageMetadata.ContentAddress = *contentAddress.S
			}

			updatedItem, err := dynamodbattribute.MarshalMap(&msg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Printf(" | %#v", updatedItem)

			updateInput := &dynamodb.PutItemInput{
				TableName: aws.String(db.TableName),
				Item:      updatedItem,
			}

			_, err = db.Client.PutItem(updateInput)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		if output.LastEvaluatedKey == nil {
			log.Info("\nNo more items to scan")
			break
		}

		input.ExclusiveStartKey = output.LastEvaluatedKey
	}

}

// 	for _, item := range result.Items {
// 		message := Message{}
//
// 		err = dynamodbattribute.UnmarshalMap(item, &message)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
//
// 		// Set UserID to AuthorID
// 		message.UserID = message.AuthorID
//
// 		// Marshal the updated message
// 		updatedItem, err := dynamodbattribute.MarshalMap(message)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
//
// 		// Create the DynamoDB UpdateItem input
// 		updateInput := &dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      updatedItem,
// 		}
//
// 		_, err = svc.PutItem(updateInput)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 	}
//
// 	fmt.Println("UserID attributes updated successfully.")
// }
