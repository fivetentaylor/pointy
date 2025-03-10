package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := dynamo.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %+v", err)
	}

	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(db.TableName),
	}

	tableDescription, err := db.Client.DescribeTable(describeTableInput)
	if err != nil {
		log.Fatalf("Failed to describe table: %v", err)
	}

	// Check if GSI1PK exists in the attribute definitions
	gsi1pkExists := false
	for _, attribute := range tableDescription.Table.AttributeDefinitions {
		if *attribute.AttributeName == "GSI1PK" {
			gsi1pkExists = true
			break
		}
	}

	// If GSI1PK exists, exit early
	if gsi1pkExists {
		log.Info("GSI1PK already exists, no need to update the table or migrate data.")
		return
	}

	log.Info("Adding GSI1PK to the table")

	updateTableInput := &dynamodb.UpdateTableInput{
		TableName: aws.String(db.TableName),
		// Add new attribute definitions
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GSI1SK"),
				AttributeType: aws.String("S"),
			},
		},
		// Specify the new GSI
		GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{
			{
				Create: &dynamodb.CreateGlobalSecondaryIndexAction{
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
						ProjectionType: aws.String("ALL"), // Adjust as needed
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
		},
	}

	_, err = db.Client.UpdateTable(updateTableInput)
	if err != nil {
		log.Fatalf("Failed to update table: %v", err)
	}

	log.Info("Migration completed successfully")
}
