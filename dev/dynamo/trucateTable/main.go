package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := dynamo.NewDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	items, err := db.ScanTable()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Iterate over the results and delete each item
	for _, item := range items {
		key := map[string]*dynamodb.AttributeValue{
			"PK": {
				S: item["PK"].S,
			},
			"SK": {
				S: item["SK"].S,
			},
		}

		deleteParams := &dynamodb.DeleteItemInput{
			TableName: aws.String(db.TableName),
			Key:       key,
		}

		_, err := db.Client.DeleteItem(deleteParams)
		if err != nil {
			fmt.Print("!")
			continue
		}

		fmt.Print(".")
	}
}
