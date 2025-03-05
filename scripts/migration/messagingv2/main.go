package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

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

		log.Infof("Scanned %d items", len(output.Items)) // output.Items

		for _, item := range output.Items {
			fmt.Printf("|")
			fmt.Printf("%s", *item["PK"].S)
			// newItem, err := updateItem(item)
			// if err != nil {
			// 	log.Fatalf("failed to update item: %+v", err)
			// }
			//
			// if newItem == nil {
			// 	if !strings.HasPrefix(*item["SK"].S, "log#") {
			// 		fmt.Printf("?")
			// 		continue
			// 	}
			// }
			//
			// if newItem != nil {
			// 	// save new item
			// 	_, err = db.Client.PutItem(&dynamodb.PutItemInput{
			// 		TableName: &db.TableName,
			// 		Item:      newItem,
			// 	})
			// 	if err != nil {
			// 		log.Fatalf("failed to update item: %+v", err)
			// 	}
			// 	fmt.Printf(".")
			// }
			//
			// // delete existing item
			// _, err = db.Client.DeleteItem(&dynamodb.DeleteItemInput{
			// 	TableName: &db.TableName,
			// 	Key: map[string]*dynamodb.AttributeValue{
			// 		"PK": {
			// 			S: item["PK"].S,
			// 		},
			// 		"SK": {
			// 			S: item["SK"].S,
			// 		},
			// 	},
			// })
			// if err != nil {
			// 	log.Fatalf("failed to delete item: %+v", err)
			// }
			// fmt.Printf("!")
		}

		if output.LastEvaluatedKey == nil {
			log.Info("\nNo more items to scan")
			break
		}

		input.ExclusiveStartKey = output.LastEvaluatedKey
	}

	log.Info("Done 🤙")
}

func printItem(item map[string]*dynamodb.AttributeValue) {
	attrs := map[string]any{}
	for k, v := range item {
		if k == "attachments" || k == "content" {
			attrs[k] = fmt.Sprintf("len(%d)", len(*v.S))
			continue
		}
		if v.S != nil {
			attrs[k] = *v.S
			continue
		}
		if v.N != nil {
			attrs[k] = *v.N
			continue
		}
		attrs[k] = v
	}

	log.Infof("%+v", attrs)
}

func updateItem(item map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	pk := *item["PK"].S
	sk := *item["SK"].S
	sk1 := *item["SK1"].S

	out := map[string]*dynamodb.AttributeValue{}

	// message
	if strings.HasPrefix(pk, "thrd#") && strings.HasPrefix(sk, "msg#") {
		// top level message
		if strings.HasPrefix(sk, "msg#top") {
			newPK := strings.Replace(pk, "thrd", "chan", 1)
			out["PK"] = &dynamodb.AttributeValue{S: aws.String(newPK)}

			createdAt := strings.Split(*item["SK"].S, "#")[2]
			out["SK1"] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprintf("msg#main@%s", createdAt))}

			out["SK"] = item["SK1"]
		}

		channelID := strings.Split(pk, "#")[1]

		if strings.HasPrefix(sk, "msg#reply_") {
			containerID := strings.Split(sk, "#")[1]
			// remove "reply_" prefix
			if strings.HasPrefix(containerID, "reply_") {
				containerID = containerID[6:]
			}

			out["PK"] = &dynamodb.AttributeValue{S: aws.String(
				fmt.Sprintf("msg#%s", containerID),
			)}
			out["SK"] = item["SK1"]

			createdAt := strings.Split(*item["SK"].S, "#")[2]
			out["SK1"] = &dynamodb.AttributeValue{S: aws.String(
				fmt.Sprintf("msg#main@%s", createdAt),
			)}

			out["parentContainerID"] = &dynamodb.AttributeValue{S: aws.String(
				fmt.Sprintf("chan#%s", channelID),
			)}
		}

		out["channelID"] = &dynamodb.AttributeValue{S: aws.String(channelID)}
		out["authorID"] = item["authorID"]
		out["attachments"] = item["attachments"]
		out["content"] = item["content"]
		out["replyCount"] = item["replyCount"]
		out["lifecycleStage"] = item["status"]

		return out, nil
	}

	// thread -> Channel
	if strings.HasPrefix(pk, "doc#") && strings.HasPrefix(sk, "thrd#") {
		out["PK"] = item["PK"]
		out["SK"] = &dynamodb.AttributeValue{S: aws.String(strings.Replace(sk, "thrd", "chan", 1))}
		out["SK1"] = &dynamodb.AttributeValue{S: aws.String(strings.Replace(sk1, "thrd", "chan", 1))}
		out["userIDs"] = item["userIDs"]

		return out, nil
	}

	return nil, nil
}
