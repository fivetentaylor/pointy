package dynamo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func printWorkdir() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Current working directory:", dir)
}

func TestMain(m *testing.M) {
	printWorkdir()
	if err := godotenv.Load("../../../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file")
	}

	db, err := NewDB()
	if err != nil {
		log.Fatalf("failed to create db: %+v", err)
	}
	err = db.CreateTestTable()
	if err != nil {
		log.Fatalf("failed to create table: %+v", err)
	}

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

func TestPutGetItem(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	id := uuid.New().String()

	item := map[string]*dynamodb.AttributeValue{
		"PK":      {S: aws.String("doc#someDocId")},
		"SK":      {S: aws.String(fmt.Sprintf("log#%s", id))},
		"TeskKey": {S: aws.String("hello world")},
	}
	err = db.PutUniqueItem(item)
	assert.NoError(t, err)

	gotItem, err := db.GetItem("doc#someDocId", fmt.Sprintf("log#%s", id))
	assert.NoError(t, err)

	jsonBytes, err := ItemToJson(gotItem)
	assert.NoError(t, err)

	log.Print(string(jsonBytes))
}
